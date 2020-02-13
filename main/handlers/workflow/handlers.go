package workflow

import (
	"log"
	"net/http"

	"github.com/ProxeusApp/proxeus-core/storage/portable"

	"github.com/ProxeusApp/proxeus-core/storage"
	"github.com/ProxeusApp/proxeus-core/sys/workflow"

	"github.com/labstack/echo"

	"github.com/ProxeusApp/proxeus-core/main/handlers/api"

	"github.com/ProxeusApp/proxeus-core/main/handlers/externalnode"
	"github.com/ProxeusApp/proxeus-core/main/handlers/helpers"
	"github.com/ProxeusApp/proxeus-core/main/www"
	"github.com/ProxeusApp/proxeus-core/sys/model"
)

func ExportWorkflow(e echo.Context) error {
	c := e.(*www.Context)
	sess := c.Session(false)
	if sess == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	var id []string
	if c.QueryParam("id") != "" {
		id = []string{c.QueryParam("id")}
	} else if c.QueryParam("contains") != "" {
		items, _ := c.System().DB.Workflow.List(sess, c.QueryParam("contains"), storage.Options{Limit: 1000})
		if len(items) > 0 {
			id = make([]string, len(items))
			for i, item := range items {
				id[i] = item.ID
			}
		}
	}
	return api.Export(sess, []portable.EntityType{portable.Workflow}, c, id...)
}

func GetHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	sess := c.Session(true)
	if sess != nil {
		item, err := c.System().DB.Workflow.Get(sess, ID)
		if err == nil {
			workflowOwner, err := c.System().DB.User.Get(sess, item.Owner)
			if err != nil {
				return c.NoContent(http.StatusNotFound)
			}
			item.OwnerEthAddress = workflowOwner.EthereumAddr
			return c.JSON(http.StatusOK, item)
		}
	}
	return c.NoContent(http.StatusNotFound)
}

func UpdateHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.QueryParam("id")
	publish := false
	if _, ok := c.QueryParams()["publish"]; ok {
		publish = true
	}
	sess := c.Session(false)
	if sess != nil {
		item := &model.WorkflowItem{}
		if err := c.Bind(&item); err != nil {
			//WorkflowItem.Price int overflow will return in bind error
			log.Println("[workflowHandler][UpdateHandler] ", err.Error())
			return c.String(http.StatusBadRequest, "unable to bind request")
		}
		item.ID = ID
		if item.Price < 0 {
			return c.String(http.StatusBadRequest, "price should be 0 or higher")
		}

		user, err := c.System().DB.User.Get(sess, sess.UserID())
		if err != nil && user == nil {
			return c.String(http.StatusBadRequest, "unable to get user")
		}
		if item.Price > 0 && user.EthereumAddr == "" {
			return c.String(http.StatusBadRequest, "can not set price without eth addr")
		}

		if publish {
			errs := map[string]interface{}{}
			collectError := func(err error, node *workflow.Node) {
				errs[node.ID] = struct {
					Error string
					Item  interface{}
				}{Error: err.Error(), Item: node}
			}
			//loop recursively and change permissions on all children
			item.LoopNodes(nil, func(l *workflow.Looper, node *workflow.Node) bool {
				if node.Type == "form" {
					it, er := c.System().DB.Form.Get(sess, node.ID)
					if er != nil {
						collectError(er, node)
						return true //continue
					}
					if !it.Published {
						it.Published = true
						er = c.System().DB.Form.Put(sess, it)
						if er != nil {
							collectError(er, node)
						}
					}
				} else if node.Type == "template" {
					it, er := c.System().DB.Template.Get(sess, node.ID)
					if er != nil {
						collectError(er, node)
						return true //continue
					}
					if !it.Published {
						it.Published = true
						er = c.System().DB.Template.Put(sess, it)
						if er != nil {
							collectError(er, node)
						}
					}
				} else if node.Type == "workflow" { // deep dive...
					it, er := c.System().DB.Workflow.Get(sess, node.ID)
					if er != nil {
						collectError(er, node)
						return true //continue
					}
					if !it.Published {
						it.Published = true
						er = c.System().DB.Workflow.Put(sess, it)
						if er != nil {
							collectError(er, node)
						}
					}
					it.LoopNodes(l, nil)
				}
				return true //continue
			})
			if len(errs) > 0 {
				return c.JSON(http.StatusMultiStatus, errs)
			}
		}

		err = c.System().DB.Workflow.Put(sess, item)
		if err != nil {
			if err == model.ErrAuthorityMissing {
				return c.NoContent(http.StatusUnauthorized)
			}
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusOK, item)
	}
	return c.NoContent(http.StatusBadRequest)
}

func DeleteHandler(e echo.Context) error {
	c := e.(*www.Context)
	ID := c.Param("ID")
	sess := c.Session(false)
	if sess != nil {
		err := c.System().DB.Workflow.Delete(sess, ID)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
	return c.NoContent(http.StatusBadRequest)
}

func ListPublishedHandler(e echo.Context) error {
	return listHandler(e.(*www.Context), true)
}

func ListHandler(e echo.Context) error {
	return listHandler(e.(*www.Context), false)
}

func listHandler(c *www.Context, publishedOnly bool) error {
	var sess model.Auth
	if s := c.Session(false); s != nil {
		sess = s
	}
	contains := c.QueryParam("c")
	settings := helpers.RequestOptions(c)
	var dat []*model.WorkflowItem
	var err error

	if publishedOnly {
		dat, err = c.System().DB.Workflow.ListPublished(sess, contains, settings)
	} else {
		dat, err = c.System().DB.Workflow.List(sess, contains, settings)
	}

	if err != nil {
		if err == model.ErrAuthorityMissing {
			log.Println("Can't list workflows: " + err.Error())
			return c.NoContent(http.StatusUnauthorized)
		}
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, dat)
}

func ListCustomNodeHandler(e echo.Context) error {
	c := e.(*www.Context)
	nodeType := c.Param("type")
	sess := c.Session(false)
	if sess != nil {
		externalnode.ProbeExternalNodes(c.System())
		dat := externalnode.List(c, nodeType)
		if dat != nil {
			return c.JSON(http.StatusOK, dat)
		}
	}
	return c.NoContent(http.StatusNotFound)
}
