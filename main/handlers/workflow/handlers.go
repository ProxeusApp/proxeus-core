package workflow

import (
	"log"
	"net/http"

	"github.com/ProxeusApp/proxeus-core/service"

	"github.com/ProxeusApp/proxeus-core/storage/portable"

	"github.com/labstack/echo"

	extNode "github.com/ProxeusApp/proxeus-core/externalnode"

	"github.com/ProxeusApp/proxeus-core/storage"

	"github.com/ProxeusApp/proxeus-core/main/handlers/api"

	"github.com/ProxeusApp/proxeus-core/main/handlers/helpers"
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

var (
	workflowService service.WorkflowService
	userService     service.UserService
	nodeService     service.NodeService
)

func Init(workflowS service.WorkflowService, userS service.UserService, nodeS service.NodeService) {
	workflowService = workflowS
	userService = userS
	nodeService = nodeS
}

func ExportWorkflow(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	var (
		id  []string
		err error
	)
	if c.QueryParam("id") != "" {
		id = []string{c.QueryParam("id")}
	} else if c.QueryParam("contains") != "" {
		id, err = workflowService.ListIds(sess, c.QueryParam("contains"), storage.Options{Limit: 1000})
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
	}
	return api.Export(sess, []portable.EntityType{portable.Workflow}, c, id...)
}

func GetHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	sess := c.Session(true)
	if sess == nil {
		return c.NoContent(http.StatusNotFound)
	}

	workflow, err := workflowService.GetAndPopulateOwner(sess, ID)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, workflow)
}

func UpdateHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.QueryParam("id")
	publish := false
	if _, ok := c.QueryParams()["publish"]; ok {
		publish = true
	}
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusBadRequest)
	}

	workflowItem := &model.WorkflowItem{}
	if err := c.Bind(&workflowItem); err != nil {
		//WorkflowItem.Price int overflow will return in bind error
		log.Println("[workflowHandler][UpdateHandler] ", err.Error())
		return c.String(http.StatusBadRequest, "unable to bind request")
	}

	workflowItem.ID = ID
	if workflowItem.Price < 0 {
		return c.String(http.StatusBadRequest, "price should be 0 or higher")
	}

	user, err := userService.GetUser(sess)
	if err != nil && user == nil {
		return c.String(http.StatusBadRequest, "unable to get user")
	}
	if workflowItem.Price > 0 && user.EthereumAddr == "" {
		return c.String(http.StatusBadRequest, "can not set price without eth addr")
	}

	if publish {
		errs := workflowService.Publish(sess, workflowItem)
		if len(errs) > 0 {
			return c.JSON(http.StatusMultiStatus, errs)
		}
	}

	instantiateExternalNode(c, sess, item)

	err = workflowService.Put(sess, workflowItem)
	if err != nil {
		if err == model.ErrAuthorityMissing {
			return c.NoContent(http.StatusUnauthorized)
		}
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, workflowItem)
}

func instantiateExternalNode(c *www.Context, auth model.Auth, item *model.WorkflowItem) {
	item.LoopNodes(nil, func(l *workflow.Looper, node *workflow.Node) bool {
		if node.Type == "externalNode" {
			_, er := c.System().DB.Workflow.QueryFromInstanceID(auth, node.ID)
			if er == nil {
				log.Println("UpdateHandler externalNode instance exists, will not create new instance")
				return true
			}

			newExternalNode, err := c.System().DB.Workflow.NodeByName(auth, node.Name)
			if err != nil {
				log.Println("UpdateHandler instantiateExternalNode NodeByName err: ", err.Error())
				return true
			}

			i := &extNode.ExternalNodeInstance{
				ID:       node.ID,
				NodeName: newExternalNode.Name,
			}

			//PutExternalNodeInstance
			er = c.System().DB.Workflow.PutExternalNodeInstance(auth, i)
			if er != nil {
				log.Println("UpdateHandler externalNode PutExternalNodeInstance error: ", er.Error())
				return true //continue
			}
		}
		return true
	})
}

func DeleteHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	sess := c.Session(false)
	if sess != nil {
		err := workflowService.Delete(sess, ID)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
	return c.NoContent(http.StatusBadRequest)
}

func ListPublishedHandler(e echo.Context) error {
	return listHandler(e.(*www.Context), true, e.QueryParam("c"), helpers.RequestOptions(e))
}

func ListHandler(e echo.Context) error {
	return listHandler(e.(*www.Context), false, e.QueryParam("c"), helpers.RequestOptions(e))
}

func listHandler(c *www.Context, published bool, contains string, settings storage.Options) error {
	var (
		err           error
		workflowItems []*model.WorkflowItem
	)
	if published {
		workflowItems, err = workflowService.ListPublished(c.Session(false), contains, settings)
	} else {
		workflowItems, err = workflowService.List(c.Session(false), contains, settings)
	}
	if err != nil {
		if err == model.ErrAuthorityMissing {
			log.Println("Can't list workflows: " + err.Error())
			return c.NoContent(http.StatusUnauthorized)
		}
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, workflowItems)
}

func ListCustomNodeHandler(e echo.Context) error {
	c := e.(*www.Context)
	nodeType := c.Param("type")
	sess := c.Session(false)
	if sess != nil {
		go externalnode.ProbeExternalNodes(c.System())
		dat := externalnode.List(c, nodeType)
		if dat != nil {
			return c.JSON(http.StatusOK, dat)
		}
	}

	nodeService.ProbeExternalNodes()
	nodes := nodeService.List(nodeType)
	if nodes == nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, nodes)
}
