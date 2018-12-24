package daos

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/Proofsuite/amp-matching-engine/types"
	"github.com/globalsign/mgo/dbtest"
)

var server dbtest.DBServer

func init() {
	temp, _ := ioutil.TempDir("", "test")
	server.SetPath(temp)

	session := server.Session()
	db = &Database{session}
}

func TestWalletDao(t *testing.T) {
	key := "7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712660"
	w := types.NewWalletFromPrivateKey(key)
	dao := NewWalletDao()

	err := dao.Create(w)
	if err != nil {
		t.Errorf("Could not create wallet object")
	}

	all, err := dao.GetAll()
	if err != nil {
		t.Errorf("Could not get wallets: %v", err)
	}

	if !reflect.DeepEqual(w, &all[0]) {
		t.Errorf("Could not retrieve correct wallets:\n Expected: %v\n, Got: %v\n", w, &all[0])
	}

	byId, err := dao.GetByID(w.ID)
	if err != nil {
		t.Errorf("Could not get wallet by ID: %v", err)
	}

	if !reflect.DeepEqual(w, byId) {
		t.Errorf("Could not correct walley by ID:\n Expected: %v\n, Got: %v\n", w, byId)
	}

	byAddress, err := dao.GetByAddress(w.Address)
	if err != nil {
		t.Errorf("Could not get wallet by address: %v", err)
	}

	if !reflect.DeepEqual(w, byAddress) {
		t.Errorf("Could not get correct wallet by address:\n Expected: %v\n, Got: %v\n", w, byAddress)
	}
}

func TestDefaultAdminWallet(t *testing.T) {
	key := "7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712660"
	w := types.NewWalletFromPrivateKey(key)
	w.Admin = true
	dao := NewWalletDao()

	err := dao.Create(w)
	if err != nil {
		t.Errorf("Could not create wallet object")
	}

	wallet, err := dao.GetDefaultAdminWallet()
	if err != nil {
		t.Errorf("Could not get default admin wallet")
	}

	if !reflect.DeepEqual(w, wallet) {
		t.Errorf("Could not get correct admin wallet:\n Expected: %v\n, Got: %v\n", w, wallet)
	}
}
