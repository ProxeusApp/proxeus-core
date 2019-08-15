package workflow

import (
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"git.proxeus.com/core/central/sys/workflow"

	"github.com/labstack/echo"

	"git.proxeus.com/core/central/main/handlers/api"

	"git.proxeus.com/core/central/main/customNode"
	"git.proxeus.com/core/central/main/helpers"
	"git.proxeus.com/core/central/main/www"
	"git.proxeus.com/core/central/sys/db/storm"
	"git.proxeus.com/core/central/sys/model"
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
		items, _ := c.System().DB.Workflow.List(sess, c.QueryParam("contains"), map[string]interface{}{"limit": 1000})
		if len(items) > 0 {
			id = make([]string, len(items))
			for i, item := range items {
				id[i] = item.ID
			}
		}
	}
	return api.Export(sess, []storm.ImexIF{c.System().DB.Workflow}, c, id...)
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
			usrAuth := &model.User{Role: model.USER}
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
					if !it.IsReadGrantedFor(usrAuth) {
						it.Permissions = item.Permissions
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
					if !it.IsReadGrantedFor(usrAuth) {
						it.Permissions = item.Permissions
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
					if !it.IsReadGrantedFor(usrAuth) {
						it.Permissions = item.Permissions
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

func GetWorkflowPayment(e echo.Context) error {
	c := e.(*www.Context)
	txHash := c.QueryParam("txHash")
	workflowId := c.Param("ID")

	var (
		workflowPaymentItem *model.WorkflowPaymentItem
		err                 error
	)
	if txHash == "" {
		sess := c.Session(false)
		user, err := c.System().DB.User.Get(sess, sess.UserID())
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
		workflowPaymentItem, err = c.System().DB.WorkflowPaymentsDB.GetByWorkflowIdAndFromEthAddress(workflowId, user.EthereumAddr)
		if err != nil {
			if err.Error() == "not found" {
				return c.NoContent(http.StatusNotFound)
			}
			return c.NoContent(http.StatusBadRequest)
		}
		err = checkPayment(c, workflowId, workflowPaymentItem)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
	} else {
		workflowPaymentItem, err = c.System().DB.WorkflowPaymentsDB.GetByTxHash(txHash)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}
	}

	log.Println("[workflowHandler][GetWorkflowPayment]", workflowPaymentItem.Hash)

	return c.JSON(http.StatusOK, workflowPaymentItem)
}

func AddWorkflowPayment(e echo.Context) error {
	c := e.(*www.Context)
	txHash := c.Param("txHash")
	workflowId := c.Param("ID")

	workflowPaymentItem, err := c.System().DB.WorkflowPaymentsDB.GetByTxHash(txHash)
	if err != nil || workflowPaymentItem.WorkflowID != "" {
		return c.NoContent(http.StatusBadRequest)
	}

	err = checkPayment(c, workflowId, workflowPaymentItem)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	workflowPaymentItem.WorkflowID = workflowId

	err = c.System().DB.WorkflowPaymentsDB.Add(workflowPaymentItem)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)
}

var errPaymentFailed = errors.New("failed to validate payment")

func checkPayment(c *www.Context, workflowId string, workflowPaymentItem *model.WorkflowPaymentItem) error {
	sess := c.Session(false)
	if sess == nil {
		return errPaymentFailed
	}
	workflow, err := c.System().DB.Workflow.Get(sess, workflowId)
	if err != nil {
		return err
	}

	if workflowPaymentItem.Xes != workflow.Price {
		return errPaymentFailed
	}

	payer, err := c.System().DB.User.Get(sess, sess.UserID())
	if err != nil || payer == nil {
		return errPaymentFailed
	}

	if payer.EthereumAddr == "" {
		return errPaymentFailed
	}

	if !strings.EqualFold(workflowPaymentItem.From, payer.EthereumAddr) {
		return errPaymentFailed
	}

	workflowOwner, err := c.System().DB.User.Get(sess, workflow.Owner)
	if err != nil {
		return errPaymentFailed
	}

	if !strings.EqualFold(workflowPaymentItem.To, workflowOwner.EthereumAddr) {
		return errPaymentFailed
	}

	return nil
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

func ListHandler(e echo.Context) error {
	c := e.(*www.Context)
	contains := c.QueryParam("c")
	sess := c.Session(false)
	if sess != nil {
		settings := helpers.ReadReqSettings(c)
		dat, err := c.System().DB.Workflow.List(sess, contains, settings)
		if err != nil {
			log.Println("Can't list workflows: " + err.Error())
			if err == model.ErrAuthorityMissing {
				return c.NoContent(http.StatusUnauthorized)
			}
			return c.NoContent(http.StatusNotFound)
		}
		return c.JSON(http.StatusOK, dat)
	}
	return c.NoContent(http.StatusUnauthorized)
}

func ListCustomNodeHandler(e echo.Context) error {
	c := e.(*www.Context)
	nodeType := c.Param("type")
	sess := c.Session(false)
	if sess != nil {
		dat := customNode.List(nodeType)
		if dat != nil {
			return c.JSON(http.StatusOK, []interface{}{dat})
		}
	}
	return c.NoContent(http.StatusNotFound)
}
