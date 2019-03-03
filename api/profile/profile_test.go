package profile_test

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/machmum/gorest/api/profile"
	"github.com/machmum/gorest/config"
	"github.com/machmum/gorest/utl/mock/mockdb"
	"github.com/machmum/gorest/utl/model"
	"github.com/machmum/gorest/utl/model/postgresql"
	"github.com/machmum/gorest/utl/platform/postgresql"
	"testing"
)

func TestCategory(t *testing.T) {
	// initiate new config
	cfg := &config.RDB{
		Lifetime: config.RDBLifetime{
			Token: 60,
			Apps:  60,
		},
		Prefix: config.RDBPrefix{
			Access:  "ce_access_",
			Refresh: "ce_refresh_",
			Apps:    "ce_apps",
		},
	}

	pss := make(gorestdb.ProfileSimpleSlice, 1)
	prs := make(gorest.ProductSlice, 1)
	result := make(gorest.ProfileSlice, 1)

	for i := range pss {
		pss[i] = gorestdb.ProfileSimple{
			ID:             1,
			FirstName:      "user",
			LastName:       "user",
			Username:       "username",
			ProductID:      1,
			ProductName:    "product ck1-a",
			NormalPrice:    "1000",
			SalePrice:      "100",
			ProductImageID: 1,
			ProductImage:   "https://tinyjpg.com/images/social/website.jpg",
		}

		prs[i] = gorest.Product{
			ID:   1,
			Name: "product ck1-a",
			Price: &gorest.Price{
				Normal: "1000",
				Sale:   "100",
			},
			Image: &gorest.ImageProduct{
				Thumbnail: &gorest.ImageRes{
					URL: "https://tinyjpg.com/images/social/website.jpg",
					Size: &gorest.WHReq{
						Width:  1020,
						Height: 510,
					},
				},
			},
		}

		result[i] = gorest.Profile{
			ID:        pss[0].ID,
			FirstName: pss[0].FirstName,
			LastName:  pss[0].LastName,
			Username:  pss[0].Username,
			Products:  prs,
		}

	}

	profileDB := &mockdb.Profile{
		FindByCategoryReturnSimpleFn: func(db *gorm.DB, category string) (slice gorestdb.ProfileSimpleSlice, err error) {
			return pss, nil
		},
	}

	cases := []struct {
		name     string
		platform *profile.Platform
		internal *profile.Internal
		result   gorest.ProfileSlice
		err      error
	}{
		{
			name: "Failed to get profile",
			err:  postgresql.ErrNotFoundProfile,
			platform: &profile.Platform{
				Profile: &mockdb.Profile{
					FindByCategoryReturnSimpleFn: func(db *gorm.DB, category string) (slice gorestdb.ProfileSimpleSlice, err error) {
						return nil, postgresql.ErrNotFoundProfile
					},
				},
			},
			internal: &profile.Internal{
				Product: nil,
			},
		},
		{
			name: "Success",
			platform: &profile.Platform{
				Profile: profileDB,
			},
			internal: &profile.Internal{
				Product: nil,
			},
			result: result,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := profile.New(nil, nil, cfg, tt.platform, tt.internal)
			r, err := s.Profile(nil, "user")

			if r != nil {
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
// 		profileDB := &mockdb.Profile{
// 			FindByCategoryReturnSimpleFn: func(db *gorm.DB, url string) (gorestdb.ProfileSimpleSlice, error) {
// 				pss := make(gorestdb.ProfileSimpleSlice, 1)
// 				for i := range pss {
// 					pss[i] = gorestdb.ProfileSimple{
// 						ID:             18,
// 						FirstName:      "cooker 1",
// 						LastName:       "",
// 						Username:       "cooker 1",
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
// 		s := profile.New(nil, nil, cfg, profileDB)
// 		_, _ = s.Profile(nil, "cook", []int{})
// 	}
// }
