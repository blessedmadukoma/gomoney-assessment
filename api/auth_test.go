package api

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	db "github.com/blessedmadukoma/gomoney-assessment/db/models"
// 	"github.com/blessedmadukoma/gomoney-assessment/utils"
// 	"github.com/gin-gonic/gin"

// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"

// 	"github.com/stretchr/testify/require"
// )

// func TestRegisterAPI(t *testing.T) {
// 	user, password := randomUser(t)

// 	testCases := []struct {
// 		name          string
// 		body          gin.H
// 		buildStubs    func(collections map[string]*mongo.Collection)
// 		checkResponse func(recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "OK",
// 			body: gin.H{
// 				"firstname": user.FirstName,
// 				"lastname":  user.LastName,
// 				"password":  password,
// 				"role":      "fan",
// 				"email":     user.Email,
// 			},
// 			buildStubs: func(collections map[string]*mongo.Collection) {
// 				collections["users"].(*MockCollection).InsertOneFunc = func(interface{}) (*mongo.InsertOneResult, error) {
// 					return &mongo.InsertOneResult{}, nil
// 				}
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusCreated, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "InternalError",
// 			body: gin.H{
// 				"firstname": user.FirstName,
// 				"lastname":  user.LastName,
// 				"password":  password,
// 				"role":      "fan",
// 				"email":     user.Email,
// 			},
// 			buildStubs: func(collections map[string]*mongo.Collection) {
// 				collections["users"].(*MockCollection).InsertOneFunc = func(interface{}) (*mongo.InsertOneResult, error) {
// 					return nil, fmt.Errorf("error inserting user")
// 				}
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "DuplicateUsername",
// 			body: gin.H{
// 				"firstname": user.FirstName,
// 				"lastname":  user.LastName,
// 				"password":  password,
// 				"role":      "fan",
// 				"email":     user.Email,
// 			},
// 			buildStubs: func(collections map[string]*mongo.Collection) {
// 				collections["users"].(*MockCollection).InsertOneFunc = func(interface{}) (*mongo.InsertOneResult, error) {
// 					return nil, mongo.WriteException{WriteErrors: []mongo.WriteError{{Code: 11000}}}
// 				}
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 			},
// 		},
// 		// Add more test cases as needed
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			collections := make(map[string]*mongo.Collection)
// 			collections["users"] = &MockCollection{}

// 			tc.buildStubs(collections)

// 			server := newTestServer(t, collections)
// 			recorder := httptest.NewRecorder()

// 			data, err := json.Marshal(tc.body)
// 			require.NoError(t, err)

// 			url := "/api/auth/register"
// 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// 			require.NoError(t, err)

// 			server.router.ServeHTTP(recorder, request)
// 			tc.checkResponse(recorder)
// 		})
// 	}
// }

// // // type eqCreateUserParamsMatcher struct {
// // // 	arg      db.RegisterUserParams
// // // 	password string
// // // }

// // // func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
// // // 	arg, ok := x.(db.RegisterUserParams)
// // // 	if !ok {
// // // 		return false
// // // 	}

// // // 	err := utils.VerifyPassword(e.password, arg.Password)
// // // 	if err != nil {
// // // 		return false
// // // 	}

// // // 	e.arg.Password = arg.Password
// // // 	return reflect.DeepEqual(e.arg, arg)
// // // }

// // // func (e eqCreateUserParamsMatcher) String() string {
// // // 	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
// // // }

// // // func EqCreateUserParams(arg db.RegisterUserParams, password string) gomock.Matcher {
// // // 	return eqCreateUserParamsMatcher{arg, password}
// // // }

// // func TestRegisterAPI(t *testing.T) {
// // 	user, password := randomUser(t)

// // 	testCases := []struct {
// // 		name          string
// // 		body          gin.H
// // 		buildStubs    func(collections map[string]*mongo.Collection)
// // 		checkResponse func(recoder *httptest.ResponseRecorder)
// // 	}{
// // 		{
// // 			name: "OK",
// // 			body: gin.H{
// // 				"firstname": user.FirstName,
// // 				"lastname":  user.LastName,
// // 				"password":  password,
// // 				"role":      user.Role,
// // 				"email":     user.Email,
// // 			},
// // 			buildStubs: func(collections map[string]*mongo.Collection) {
// // 				collections["users"].InsertOne(nil, &user)
// // 			},
// // 			// 	arg := db.CreateUserParams{
// // 			// 		Username: user.Username,
// // 			// 		FullName: user.FullName,
// // 			// 		Email:    user.Email,
// // 			// 	}
// // 			// 	store.EXPECT().
// // 			// 		CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
// // 			// 		Times(1).
// // 			// 		Return(user, nil)
// // 			// },
// // 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// // 				require.Equal(t, http.StatusOK, recorder.Code)
// // 				requireBodyMatchUser(t, recorder.Body, user)
// // 			},
// // 		},
// // 		{
// // 			name: "InternalError",
// // 			body: gin.H{
// // 				"firstname": user.FirstName,
// // 				"lastname":  user.LastName,
// // 				"password":  password,
// // 				"role":      user.Role,
// // 				"email":     user.Email,
// // 			},
// // 			buildStubs: func(collections map[string]*mongo.Collection) {
// // 				collections["users"].InsertOne(nil, &user)
// // 			},
// // 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// // 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// // 			},
// // 		},
// // 		{
// // 			name: "DuplicateUsername",
// // 			body: gin.H{
// // 				"firstname": user.FirstName,
// // 				"lastname":  user.LastName,
// // 				"password":  password,
// // 				"role":      user.Role,
// // 				"email":     user.Email,
// // 			},
// // 			buildStubs: func(collections map[string]*mongo.Collection) {
// // 				collections["users"].InsertOne(nil, &user)
// // 			},
// // 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// // 				require.Equal(t, http.StatusForbidden, recorder.Code)
// // 			},
// // 		},
// // 		{
// // 			name: "InvalidFirstName",
// // 			body: gin.H{
// // 				"firstname": "invalid-user#1",
// // 				"lastname":  user.LastName,
// // 				"password":  password,
// // 				"role":      user.Role,
// // 				"email":     user.Email,
// // 			},
// // 			// buildStubs: func(store *mockdb.MockStore) {
// // 			// 	store.EXPECT().
// // 			// 		CreateUser(gomock.Any(), gomock.Any()).
// // 			// 		Times(0)
// // 			// },
// // 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// // 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// // 			},
// // 		},
// // 		{
// // 			name: "InvalidEmail",
// // 			body: gin.H{
// // 				"firstname": user.FirstName,
// // 				"lastname":  user.LastName,
// // 				"password":  password,
// // 				"role":      user.Role,
// // 				"email":     "invalid-email",
// // 			},
// // 			buildStubs: func(collections map[string]*mongo.Collection) {
// // 				collections["users"].InsertOne(nil, &user)
// // 			},
// // 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// // 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// // 			},
// // 		},
// // 		{
// // 			name: "TooShortPassword",
// // 			body: gin.H{
// // 				"firstname": user.FirstName,
// // 				"lastname":  user.LastName,
// // 				"password":  "123",
// // 				"role":      user.Role,
// // 				"email":     user.Email,
// // 			},
// // 			buildStubs: func(collections map[string]*mongo.Collection) {
// // 				collections["users"].InsertOne(nil, &user)
// // 			},
// // 			// buildStubs: func(store *mockdb.MockStore) {
// // 			// 	store.EXPECT().
// // 			// 		CreateUser(gomock.Any(), gomock.Any()).
// // 			// 		Times(0)
// // 			// },
// // 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// // 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// // 			},
// // 		},
// // 	}

// // 	for i := range testCases {
// // 		tc := testCases[i]

// // 		t.Run(tc.name, func(t *testing.T) {
// // 			// ctrl := gomock.NewController(t)
// // 			// defer ctrl.Finish()

// // 			// store := mockdb.NewMockStore(ctrl)
// // 			collections := make(map[string]*mongo.Collection)
// // 			tc.buildStubs(collections)

// // 			server := newTestServer(t, make(map[string]*mongo.Collection))
// // 			recorder := httptest.NewRecorder()

// // 			// Marshal body data to JSON
// // 			data, err := json.Marshal(tc.body)
// // 			require.NoError(t, err)

// // 			url := "/api/users/register"
// // 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// // 			require.NoError(t, err)

// // 			server.router.ServeHTTP(recorder, request)
// // 			tc.checkResponse(recorder)
// // 		})
// // 	}
// // }

// // // func TestLoginUserAPI(t *testing.T) {
// // // 	user, password := randomUser(t)

// // // 	testCases := []struct {
// // // 		name          string
// // // 		body          gin.H
// // // 		buildStubs    func(collections map[string]*mongo.Collection)
// // // 		checkResponse func(recoder *httptest.ResponseRecorder)
// // // 	}{
// // // 		{
// // // 			name: "OK",
// // // 			body: gin.H{
// // // 				"email": user.Email,
// // // 				"password": password,
// // // 			},
// // // 			buildStubs: func(collections map[string]*mongo.Collection) {
// // // 				collections["users"].FindOne(nil, &user)
// // // 			},
// // // 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// // // 				require.Equal(t, http.StatusOK, recorder.Code)
// // // 			},
// // // 		},
// // // 		{
// // // 			name: "UserNotFound",
// // // 			body: gin.H{
// // // 				"username": "NotFound",
// // // 				"password": password,
// // // 			},
// // // 			buildStubs: func(store *mockdb.MockStore) {
// // // 				store.EXPECT().
// // // 					GetUserByUsername(gomock.Any(), gomock.Any()).
// // // 					Times(1).
// // // 					Return(db.User{}, sql.ErrNoRows)
// // // 			},
// // // 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// // // 				require.Equal(t, http.StatusNotFound, recorder.Code)
// // // 			},
// // // 		},
// // // 		{
// // // 			name: "IncorrectPassword",
// // // 			body: gin.H{
// // // 				"username": user.Username,
// // // 				"password": "incorrect",
// // // 			},
// // // 			buildStubs: func(store *mockdb.MockStore) {
// // // 				store.EXPECT().
// // // 					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
// // // 					Times(1).
// // // 					Return(user, nil)
// // // 			},
// // // 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// // // 				require.Equal(t, http.StatusUnauthorized, recorder.Code)
// // // 			},
// // // 		},
// // // 		{
// // // 			name: "InternalError",
// // // 			body: gin.H{
// // // 				"username": user.Username,
// // // 				"password": password,
// // // 			},
// // // 			buildStubs: func(store *mockdb.MockStore) {
// // // 				store.EXPECT().
// // // 					GetUserByUsername(gomock.Any(), gomock.Any()).
// // // 					Times(1).
// // // 					Return(db.User{}, sql.ErrConnDone)
// // // 			},
// // // 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// // // 				require.Equal(t, http.StatusInternalServerError, recorder.Code)
// // // 			},
// // // 		},
// // // 		{
// // // 			name: "InvalidUsername",
// // // 			body: gin.H{
// // // 				"username": "invalid-user#1",
// // // 				"password": password,
// // // 			},
// // // 			buildStubs: func(store *mockdb.MockStore) {
// // // 				store.EXPECT().
// // // 					GetUserByUsername(gomock.Any(), gomock.Any()).
// // // 					Times(0)
// // // 			},
// // // 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// // // 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// // // 			},
// // // 		},
// // // 	}

// // // 	for i := range testCases {
// // // 		tc := testCases[i]

// // // 		t.Run(tc.name, func(t *testing.T) {
// // // 			ctrl := gomock.NewController(t)
// // // 			defer ctrl.Finish()

// // // 			store := mockdb.NewMockStore(ctrl)
// // // 			tc.buildStubs(store)

// // // 			server := newTestServer(t, store)
// // // 			recorder := httptest.NewRecorder()

// // // 			// Marshal body data to JSON
// // // 			data, err := json.Marshal(tc.body)
// // // 			require.NoError(t, err)

// // // 			url := "/api/users/login"
// // // 			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
// // // 			require.NoError(t, err)

// // // 			server.router.ServeHTTP(recorder, request)
// // // 			tc.checkResponse(recorder)
// // // 		})
// // // 	}
// // // }

// func randomUser(t *testing.T) (user db.RegisterUserParams, password string) {
// 	password = utils.RandomString(6)
// 	hashedPassword, err := utils.HashPassword(password)
// 	require.NoError(t, err)

// 	user = db.RegisterUserParams{
// 		ID:        primitive.NewObjectID(),
// 		FirstName: utils.RandomString(6),
// 		LastName:  utils.RandomString(6),
// 		Email:     utils.RandomEmail(),
// 		Role:      "fan",
// 		Password:  hashedPassword,
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}
// 	return
// }

// // func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.RegisterUserParams) {
// // 	data, err := io.ReadAll(body)
// // 	require.NoError(t, err)

// // 	var gotUser db.RegisterUserParams
// // 	err = json.Unmarshal(data, &gotUser)

// // 	require.NoError(t, err)
// // 	require.Equal(t, user.FirstName, gotUser.FirstName)
// // 	require.Equal(t, user.LastName, gotUser.LastName)
// // 	require.Equal(t, user.Email, gotUser.Email)
// // 	require.Empty(t, gotUser.Password)
// // }
