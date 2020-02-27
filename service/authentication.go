package service

import (
	"github.com/ProxeusApp/proxeus-core/main/handlers/blockchain"
	"github.com/ProxeusApp/proxeus-core/storage/database/db"
	"github.com/ProxeusApp/proxeus-core/sys/model"
	"log"
)

type (
	AuthenticationService interface {
		LoginWithUsernamePassword(email, password string) (*model.User, error)
		LoginWithWallet(challenge, signature string) (bool, *model.User, error)
	}
	defaultAuthenticationService struct {
		userService     UserService
		workflowService WorkflowService
	}
)

func NewAuthenticationService(userS UserService, workflowS WorkflowService) *defaultAuthenticationService {
	return &defaultAuthenticationService{userService: userS, workflowService: workflowS}
}

// LoginWithUsernamePassword performs a login with username and password
func (me *defaultAuthenticationService) LoginWithUsernamePassword(email, password string) (*model.User, error) {
	return userDB().Login(email, password)
}

// LoginWithWallet logs the user in when the login is done with a signature (e.g. Metamask)
func (me *defaultAuthenticationService) LoginWithWallet(challenge, signature string) (bool, *model.User, error) {

	var root = &model.User{Role: model.ROOT}

	created := false
	var address string
	var err error
	address, err = blockchain.VerifySignInChallenge(challenge, signature)
	if err != nil {
		return false, nil, err
	}
	var usr *model.User
	usr, err = me.userService.GetByBCAddress(address)
	if db.NotFound(err) {
		stngs, err := settingsDB().Get()
		if err != nil {
			return false, nil, err
		}
		it := &model.User{
			EthereumAddr: address,
			Role:         model.StringToRole(stngs.DefaultRole),
		}
		it.Name = "created by ethereum account sign"
		err = userDB().Put(root, it)
		if err != nil {
			return false, nil, err
		}
		created = true
		usr, err = me.userService.GetByBCAddress(address)
		if err == nil {
			me.workflowService.CopyWorkflows(root, usr)
			if stngs.BlockchainNet == "ropsten" && stngs.AirdropEnabled == "true" {
				go func() {
					defer func() {
						if r := recover(); r != nil {
							log.Println("airdrop recover with err ", r)
						}
					}()
					blockchain.GiveTokens(address)
				}()
			}
		}

	}
	return created, usr, err
}
