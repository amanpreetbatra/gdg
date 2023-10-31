package test

import (
	"github.com/esnet/gdg/internal/service"
	"github.com/esnet/grafana-swagger-api-golang/goclient/models"
	"golang.org/x/exp/maps"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestTeamCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	if os.Getenv("TEST_TOKEN_CONFIG") == "1" {
		t.Skip("Skipping Token configuration, Team and User CRUD requires Basic Auth")
	}
	filter := service.NewTeamFilter("")
	apiClient, _, cleanup := initTest(t, nil)
	defer cleanup()
	log.Info("Exporting current user list")
	apiClient.UploadUsers(service.NewUserFilter(""))
	users := apiClient.ListUsers(service.NewUserFilter(""))
	assert.Equal(t, len(users), 2)
	log.Info("Exporting all teams")
	apiClient.UploadTeams(filter)
	log.Info("Listing all Teams")
	teamsMap := apiClient.ListTeams(filter)
	teams := maps.Keys(teamsMap)
	assert.Equal(t, len(teams), 2)
	var engineerTeam *models.TeamDTO
	var musicianTeam *models.TeamDTO
	for ndx, team := range teams {
		if team.Name == "engineers" {
			engineerTeam = teams[ndx]
		} else if team.Name == "musicians" {
			musicianTeam = teams[ndx]
		}
	}
	assert.Equal(t, engineerTeam.Name, "engineers")
	engineers := teamsMap[engineerTeam]
	assert.Equal(t, len(engineers), 2)
	assert.Equal(t, engineers[1].Login, "tux")
	assert.Equal(t, musicianTeam.Name, "musicians")
	//Import Teams
	log.Info("Importing teams")
	list := apiClient.DownloadTeams(filter)
	assert.Equal(t, len(list), len(teams))
	//CleanUp
	_, err := apiClient.DeleteTeam(filter)
	assert.Nil(t, err)
	//Remove Users
	apiClient.DeleteAllUsers(service.NewUserFilter(""))

}
