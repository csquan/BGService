package api

//func (a *ApiService) init(c *gin.Context) {
//	buf := make([]byte, 2048)
//	n, _ := c.Request.Body.Read(buf)
//	data1 := string(buf[0:n])
//	res := types.HttpRes{}
//
//	isValid := gjson.Valid(data1)
//	if isValid == false {
//		logrus.Error("Not valid json")
//		res.Code = http.StatusBadRequest
//		res.Message = "Not valid json"
//		c.SecureJSON(http.StatusBadRequest, res)
//		return
//	}
//	name := gjson.Get(data1, "name")
//	apiKey := gjson.Get(data1, "apiKey")
//	apiSecret := gjson.Get(data1, "apiSecret")
//
//	mechanismData := types.Mechanism{
//		Name:      name.String(),
//		ApiKey:    apiKey.String(),
//		ApiSecret: apiSecret.String(),
//	}
//
//	err := a.db.CommitWithSession(a.db, func(s *xorm.Session) error {
//		if err := a.db.InsertMechanism(s, &mechanismData); err != nil {
//			logrus.Errorf("insert  InsertMechanism task error:%v tasks:[%v]", err, mechanismData)
//			return err
//		}
//		return nil
//	})
//	if err != nil {
//		logrus.Error(err)
//	}
//
//	res.Code = 0
//	res.Message = err.Error()
//	res.Data = ""
//
//	c.SecureJSON(http.StatusOK, res)
//}
