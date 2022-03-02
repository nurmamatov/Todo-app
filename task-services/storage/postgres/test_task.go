package postgres

// import (
// 	"testing"
// 	"two_services/task-services/config"
// 	task "two_services/task-services/genproto"
// 	"two_services/task-services/pkg/db"
// 	"two_services/task-services/storage/repo"

// 	"github.com/stretchr/testify/suite"
// )

// type SuiteTest struct {
// 	suite.Suite
// 	Repo    repo.TaskStorageI
// 	CleanUp func()
// }

// func (suite *SuiteTest) SetupSuite() {
// 	db, _ := db.ConnectToDB(config.Load())

// 	suite.Repo = NewTaskRepo(db)

// 	suite.CleanUp = func() { db.Close() }
// }

// func (suite *SuiteTest) TestGet() {
// 	id := task.GetAndDeleteTask{Id: "716b312f-08bf-4835-8bc8-75e44e7e4ca5"}
// 	res, err := suite.Repo.Get(&id)
// 	suite.Nil(err)
// 	suite.NotNil(res)
// 	if suite.Equal(id.Id, res.Id) {
// 		suite.T().Log("error in Geting task")
// 	}
// }

// func (suite *SuiteTest) TestCreate() {
// 	new := task.CreateTaskReq{
// 		Title:      "Nimadur",
// 		AssigneeId: "716b312f-08bf-4835-8bc8-75e44e7e4ca5",
// 		Deadline:   "2022-09-09 12:23",
// 		Status:     "Activ",
// 	}
// 	res, err := suite.Repo.Create(&new)
// 	suite.Nil(err)
// 	suite.NotNil(res)
	
// }

// func (suite *SuiteTest) TearDownSuite() {
// 	suite.CleanUp()
// }

// func TestMain(T *testing.T) {
// 	suite.Run(T, new(SuiteTest))
// }
