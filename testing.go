package quirk

// testing data.
var (
	testPersonCorrect = struct {
		Username   string `quirk:"username,unique"`
		Website    string `quirk:"website"`
		AccountAge int    `quirk:"acctage"`
		Email      string `quirk:"email,unique"`
	}{
		Username:   "damienstamates",
		Website:    "github.com",
		AccountAge: 197,
		Email:      "damienstamates@gmail.com",
	}

	testPredValCorrect = []*predValDat{
		&predValDat{Predicate: "username", Value: testPersonCorrect.Username, IsUpsert: true},
		&predValDat{Predicate: "website", Value: testPersonCorrect.Website, IsUpsert: false},
		&predValDat{Predicate: "acctage", Value: testPersonCorrect.AccountAge, IsUpsert: false},
		&predValDat{Predicate: "email", Value: testPersonCorrect.Email, IsUpsert: true},
	}

	testPersonInvalid = struct {
		Username   string
		Website    string `quirk:"website"`
		AccountAge int    `quirk:"acctage"`
		Email      string `quirk:"email,unique"`
	}{
		Username:   "damienstamates",
		Website:    "github.com",
		AccountAge: 197,
		Email:      "damienstamates@gmail.com",
	}

	testPredValInvalid = []*predValDat{
		&predValDat{Predicate: "", Value: testPersonCorrect.Username, IsUpsert: false},
		&predValDat{Predicate: "website", Value: testPersonCorrect.Website, IsUpsert: false},
		&predValDat{Predicate: "acctage", Value: testPersonCorrect.AccountAge, IsUpsert: false},
		&predValDat{Predicate: "email", Value: testPersonCorrect.Email, IsUpsert: true},
	}
)
