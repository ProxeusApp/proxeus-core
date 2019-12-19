package storm

//
//func TestNew(t *testing.T) {
//	baseDir := "./testDir"
//	sdb, err := NewSettingsDB(baseDir)
//	if err != nil {
//		t.Error(err)
//	}
//	stngs := &model.Settings{}
//	stngs.DefaultRole = model.ADMIN
//	stngs.BlockchainNet = "ropsten"
//	//stngs.DatabaseDir = "./baseDir"
//	stngs.LogPath = "./log"
//	stngs.SessionExpiry = "1h"
//	//stngs.SessionManagerDir = "/managerDir"
//
//	err = sdb.Put(stngs)
//	if err != nil {
//		t.Error(err)
//	}
//	s, err := sdb.Get()
//	if err != nil {
//		t.Error(err)
//	}
//	if s == nil {
//		t.Error(s)
//	}
//	if s.DefaultRole != model.ADMIN {
//		t.Error(s)
//	}
//	err = os.RemoveAll(baseDir)
//	if err != nil {
//		t.Error(err)
//	}
//}
