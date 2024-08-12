package info

import (
	"context"
	"testing"

	"github.com/mattrmcg/equitalytics-backend/internal/db"
)

func TestInfoServiceFunctions(t *testing.T) {

	dbPool, err := db.CreateDBPool("postgres://root:123@127.0.0.1:5432/eql")
	if err != nil {
		t.Error(err)
	}

	defer db.CloseDBPool(dbPool)

	infoService := NewInfoService(dbPool)

	t.Run("should pass if GetInfoByTicker() executes correctly", func(t *testing.T) {
		companyInfo, err := infoService.GetInfoByTicker(context.Background(), "ADBE")
		if err != nil {
			t.Errorf("unable to retrieve without error: %v", err)
		}

		if companyInfo == nil {
			t.Errorf("companyInfo is returned as nil")
		}
	})
}
