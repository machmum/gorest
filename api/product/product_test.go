package product_test

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/machmum/gorest/api/product"
	"github.com/machmum/gorest/config"
	"github.com/machmum/gorest/utl/mock/mockdb"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/model/postgresql"
	"github.com/machmum/gorest/utl/platform/postgresql"
	"testing"
)

func TestProduct(t *testing.T) {
	// initiate new config
	cfg := &config.RDB{
		Lifetime: config.RDBLifetime{
			Token: 10,
			Apps:  10,
		},
		Prefix: config.RDBPrefix{
			Access:  "ce_access_",
			Refresh: "ce_refresh_",
			Apps:    "ce_apps",
		},
	}

	im := make(gorest.ImageProductSlice, 2)
	for i := range im {
		im[i] = gorest.ImageProduct{
			Thumbnail: &gorest.ImageRes{
				URL: "https://tinyjpg.com/images/social/website.jpg",
				Size: &gorest.WHReq{
					Width:  1020,
					Height: 510,
				},
			},
		}
	}

	result := gorest.Product{
		ID:   1,
		Name: "product",
		Price: &gorest.Price{
			Normal: "1000",
			Sale:   "100",
		},
		Images: im,
	}

	productDB := &mockdb.Product{
		FindByProductIDFn: func(db *gorm.DB, profileID int, pid int) (result gorestdb.ProductDetail, err error) {
			return gorestdb.ProductDetail{
				ProductID:    1,
				ProductName:  "product",
				NormalPrice:  "1000",
				SalePrice:    "100",
				ProductImage: "1+https://tinyjpg.com/images/social/website.jpg|2+https://tinyjpg.com/images/social/website.jpg",
			}, nil
		},
	}

	// logrus.Println(productDB)
	// logrus.Println(result)

	cases := []struct {
		name      string
		profileID int
		productID int
		platform  *product.Platform
		result    gorest.Product
		err       error
	}{
		{
			name:      "Failed to get profile",
			profileID: 0,
			productID: 0,
			err:       postgresql.ErrNotFoundProfile,
			platform: &product.Platform{
				Product: &mockdb.Product{
					FindByProductIDFn: func(db *gorm.DB, profileID int, pid int) (result gorestdb.ProductDetail, err error) {
						return result, postgresql.ErrNotFoundProfile
					},
				},
			},
		},
		{
			name:      "Success",
			profileID: 1,
			productID: 1,
			platform: &product.Platform{
				Product: productDB,
			},
			result: result,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := product.New(nil, nil, cfg, tt.platform)
			r, err := s.Product(nil, tt.profileID, tt.productID)

			if r.ID != 0 {
				assert.Equal(t, tt.result, r)
			}
			assert.Equal(t, tt.err, err)

		})
	}
}

// func BenchmarkCategory(b *testing.B) {
// 	b.ResetTimer()
// 	b.ReportAllocs()
// 	for n := 0; n < b.N; n++ {
// 		// initiate new config
// 		cfg := &config.RDB{
// 			Lifetime: config.RDBLifetime{
// 				Token: 10,
// 				Apps:  10,
// 			},
// 			Prefix: config.RDBPrefix{
// 				Access:  "ce_access_",
// 				Refresh: "ce_refresh_",
// 				Apps:    "ce_apps",
// 			},
// 		}
//
// 		profileDB := &mockdb.Product{
// 			FindByCategoryReturnSimpleFn: func(db *gorm.DB, url string) (gorestdb.ProfileSimpleSlice, error) {
// 				pss := make(gorestdb.ProfileSimpleSlice, 1)
// 				for i := range pss {
// 					pss[i] = gorestdb.ProfileSimple{
// 						ID:             18,
// 						FirstName:      "cooker 1",
// 						LastName:       "",
// 						Username:       "cooker 1",
// 						ProfileImage:   "http://im.berrybenka.biz/assets/cache/600/icon-mm/icon_5c2479a4b8c8e.jpg",
// 						ProductID:      1,
// 						ProductName:    "product",
// 						NormalPrice:    "1000",
// 						SalePrice:      "1000",
// 						ProductImageID: 1,
// 						ProductImage:   "http://im.berrybenka.biz/assets/cache/1125/product/zoom/274703_geslyn-sleeveless-peplum-maroon_red_M4VNB.jpg",
// 					}
// 				}
// 				return pss, nil
// 			},
// 		}
//
// 		s := product.New(nil, nil, cfg, profileDB)
// 		_, _ = s.Product(nil, "cook", []int{})
// 	}
// }

// func TestCategory(t *testing.T) {
// 	// initiate new config
// 	cfg := &config.RDB{
// 		Lifetime: config.RDBLifetime{
// 			Token: 60,
// 			Apps:  60,
// 		},
// 		Prefix: config.RDBPrefix{
// 			Access:  "ce_access_",
// 			Refresh: "ce_refresh_",
// 			Apps:    "ce_apps",
// 		},
// 	}
//
// 	// var b gorestdb.ProfileSlice
// 	b := make(gorestdb.ProfileSlice, 3)
// 	a := gorestdb.Product{
// 		ID:                     18,
// 		FirstName:              "cooker 1",
// 		LastName:               "",
// 		Password:               "$2a$08$rQ.afirxZtqnl0nRXSPfmur6zD/fh1K/.F2J3MrZ5.sOpBexCbQ4.",
// 		Username:               "cooker 1",
// 		Description:            "desc",
// 		Email:                  "cooker1@gmail.com",
// 		Phone:                  "0000",
// 		Image:                  "icon_5c2479a4b8c8e.jpg",
// 		Address:                "addr",
// 		City:                   "cty",
// 		Province:               "prov",
// 		ZipCode:                "0000",
// 		Country:                "idn",
// 		Latitude:               "111",
// 		Longitude:              "111",
// 		Verification:           9,
// 		Popular:                8,
// 		Status:                 1,
// 		StatusProfileName:      "COOKER",
// 		StatusVerificationName: "YES",
// 		StatusPopularName:      "DISABLED",
// 	}
//
// 	for i := range b {
// 		b[i] = a
// 	}
//
// 	cases := []struct {
// 		name      string
// 		profileDB *mockdb.Product
// 		scope     map[string]interface{}
// 		size      []int
// 		err       error
// 		wantError bool
// 	}{
// 		{
// 			name:  "Failed to get scope",
// 			scope: nil,
// 			err:   nil,
// 			profileDB: &mockdb.Product{
// 				FindByCategoryFn: func(db *gorm.DB, product string) (slice gorestdb.ProfileSlice, err error) {
// 					return b, nil
// 				},
// 			},
// 		},
// 	}
//
// 	for _, tt := range cases {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := product.New(nil, nil, cfg, tt.profileDB)
// 			_, err := s.Product(nil, "cook", []int{})
//
// 			assert.Equal(t, tt.err, err)
//
// 		})
// 	}
// }
