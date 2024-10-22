package steps

import (
	"context"
	"encoding/json"
	"time"

	dpMongoDriver "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/ONSdigital/dp-permissions-api/config"
	"github.com/ONSdigital/dp-permissions-api/models"

	"github.com/cucumber/godog"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	rsaJWKS = map[string]string{
		"NeKb65194Jo=": "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu1SU1LfVLPHCozMxH2Mo" +
			"4lgOEePzNm0tRgeLezV6ffAt0gunVTLw7onLRnrq0/IzW7yWR7QkrmBL7jTKEn5u" +
			"+qKhbwKfBstIs+bMY2Zkp18gnTxKLxoS2tFczGkPLPgizskuemMghRniWaoLcyeh" +
			"kd3qqGElvW/VDL5AaWTg0nLVkjRo9z+40RQzuVaE8AkAFmxZzow3x+VJYKdjykkJ" +
			"0iT9wCS0DRTXu269V264Vf/3jvredZiKRkgwlL9xNAwxXFg0x/XFw005UWVRIkdg" +
			"cKWTjpBP2dPwVZ4WWC+9aGVd+Gyn1o0CLelf4rEjGoXbAAEgAqeGUxrcIlbjXfbc" +
			"mwIDAQAB",
	}
	jwtPublisherAndAdminToken = "eyJraWQiOiJOZUtiNjUxOTRKbz0iLCJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhYWFhYWFhYS1iYmJiLWNjY2MtZGRkZC1lZWVlZWVlZWVlZWUiLCJkZXZpY2Vfa2V5IjoiYWFhYWFhYWEtYmJiYi1jY2NjLWRkZGQtZWVlZWVlZWVlZWVlIiwiY29nbml0bzpncm91cHMiOlsicm9sZS1hZG1pbiIsInJvbGUtcHVibGlzaGVyIl0sInRva2VuX3VzZSI6ImFjY2VzcyIsInNjb3BlIjoiYXdzLmNvZ25pdG8uc2lnbmluLnVzZXIuYWRtaW4iLCJhdXRoX3RpbWUiOjE1NjIxOTA1MjQsImlzcyI6Imh0dHBzOi8vY29nbml0by1pZHAudXMtd2VzdC0yLmFtYXpvbmF3cy5jb20vdXMtd2VzdC0yX2V4YW1wbGUiLCJleHAiOjQwOTc4MjUyOTIsImlhdCI6MTU2MjE5MDUyNCwianRpIjoiYWFhYWFhYWEtYmJiYi1jY2NjLWRkZGQtZWVlZWVlZWVlZWVlIiwiY2xpZW50X2lkIjoiNTdjYmlzaGs0ajI0cGFiYzEyMzQ1Njc4OTAiLCJ1c2VybmFtZSI6ImphbmVkb2VAZXhhbXBsZS5jb20ifQ.iYUem-3-HHsLaWkzN4jwcXcRV7FgKIw-_FQnY_SSnW7ZS2XH07cnZb0tDygnAyrYNd5SnHFsXsl4cXXqIG7MatQI_UawCzC07TeKvPcNpRmlTDDH-M5FBCutfmI_xAqoZbZ28oKKybrPf3KBz-DunPQ4ctF14cp1XomEEYH21GSDcn6a1xmt91PYZLxThO-zJemYMKnGNqXnm9h9F_07_8LlMIpaxdWitK1MYsXYHdaJOzi1BZX1827jGGYZBDkzQdOFCrto_XQ-Qs8SOZvkm8zIBkzL56SicuAI8VeXhn9Z7EWU4fMICmeQEPU3tSjejGLyFKfu5LXn3Q7zEbJWyw"
	jwtAdminOnlyToken         = "eyJraWQiOiJOZUtiNjUxOTRKbz0iLCJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhYWFhYWFhYS1iYmJiLWNjY2MtZGRkZC1lZWVlZWVlZWVlZWUiLCJkZXZpY2Vfa2V5IjoiYWFhYWFhYWEtYmJiYi1jY2NjLWRkZGQtZWVlZWVlZWVlZWVlIiwiY29nbml0bzpncm91cHMiOlsicm9sZS1hZG1pbiJdLCJ0b2tlbl91c2UiOiJhY2Nlc3MiLCJzY29wZSI6ImF3cy5jb2duaXRvLnNpZ25pbi51c2VyLmFkbWluIiwiYXV0aF90aW1lIjoxNTYyMTkwNTI0LCJpc3MiOiJodHRwczovL2NvZ25pdG8taWRwLnVzLXdlc3QtMi5hbWF6b25hd3MuY29tL3VzLXdlc3QtMl9leGFtcGxlIiwiZXhwIjo0MDk3ODI1MjkyLCJpYXQiOjE1NjIxOTA1MjQsImp0aSI6ImFhYWFhYWFhLWJiYmItY2NjYy1kZGRkLWVlZWVlZWVlZWVlZSIsImNsaWVudF9pZCI6IjU3Y2Jpc2hrNGoyNHBhYmMxMjM0NTY3ODkwIiwidXNlcm5hbWUiOiJqYW5lZG9lQGV4YW1wbGUuY29tIn0.Yd0HrK7pCB2UzIDDkyX1ce1G31PPAerV3DgxI7Tq7sQT4CdXMDXmeznqsLbAyfmzkDi86JPGvM_XbOGkIZ4XzvNWO2bQ4Y7B6WYYkXikADP0ge6AVu9I8B8-8n_XvxSmyC5HYwVUIPznqr7gBtOXNN1JBTcGezMCtvHrZM1f9j0WZXKi8zP6cvm4SQGqaSmSbOjZWKD1y_ZbwcWIroMO26y7vp2XxszL27vDC3FfeVk3j7gmkc5pR902wimaTocABHBieMBSi1w_u_khHfu1Ty1Idp8jghIUpd_HC_wgbS-lTxeiww19cJZysvmcbdavKYa4sP1k0aWfd5-2hetj7Q"
	jwtBasicUserOnlyToken     = "eyJraWQiOiJOZUtiNjUxOTRKbz0iLCJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhYWFhYWFhYS1iYmJiLWNjY2MtZGRkZC1lZWVlZWVlZWVlZWUiLCJkZXZpY2Vfa2V5IjoiYWFhYWFhYWEtYmJiYi1jY2NjLWRkZGQtZWVlZWVlZWVlZWVlIiwiY29nbml0bzpncm91cHMiOlsicm9sZS1iYXNpYyJdLCJ0b2tlbl91c2UiOiJhY2Nlc3MiLCJzY29wZSI6ImF3cy5jb2duaXRvLnNpZ25pbi51c2VyLmFkbWluIiwiYXV0aF90aW1lIjoxNTYyMTkwNTI0LCJpc3MiOiJodHRwczovL2NvZ25pdG8taWRwLnVzLXdlc3QtMi5hbWF6b25hd3MuY29tL3VzLXdlc3QtMl9leGFtcGxlIiwiZXhwIjo0MDk3ODI1MjkyLCJpYXQiOjE1NjIxOTA1MjQsImp0aSI6ImFhYWFhYWFhLWJiYmItY2NjYy1kZGRkLWVlZWVlZWVlZWVlZSIsImNsaWVudF9pZCI6IjU3Y2Jpc2hrNGoyNHBhYmMxMjM0NTY3ODkwIiwidXNlcm5hbWUiOiJqYW5lZG9lQGV4YW1wbGUuY29tIn0.DB9b91Xi0JNyNui8swyxJjpUm4qqyJMSjk0nxXHU51XDscL6lr2KKnU1ZwMfuYXWzZp3lyV-pVYlcDKRCATiczbOV0ZQXeVKnlePF6k18K_bUx29eVr7vBlT-0L-8Z36sCi8WPehKLdSrwDGwqp6avmSklkt9U19Sg6avwrJXsg_JFuggUKhfYewvqxw-9TkCQ9-v3BEroG5nSg6yBJgctmemTqVLnCIM7MaZPr5mUvlQPOAfzFbE80k-wXQ_eSWlNVXOsXtaVpdPfAZVAUNEKzKuPWby2y63pED-XwYb7PQM0mVgwhHpXq6RpU6ijxJmmGVUY8DYS37zAfzWZPjAw"
	jwtViewerOnlyToken        = "eyJraWQiOiJOZUtiNjUxOTRKbz0iLCJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhYWFhYWFhYS1iYmJiLWNjY2MtZGRkZC1lZWVlZWVlZWVlZWUiLCJkZXZpY2Vfa2V5IjoiYWFhYWFhYWEtYmJiYi1jY2NjLWRkZGQtZWVlZWVlZWVlZWVlIiwiY29nbml0bzpncm91cHMiOlsicm9sZS12aWV3ZXIiXSwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJhd3MuY29nbml0by5zaWduaW4udXNlci5hZG1pbiIsImF1dGhfdGltZSI6MTU2MjE5MDUyNCwiaXNzIjoiaHR0cHM6Ly9jb2duaXRvLWlkcC51cy13ZXN0LTIuYW1hem9uYXdzLmNvbS91cy13ZXN0LTJfZXhhbXBsZSIsImV4cCI6NDA5NzgyNTI5MiwiaWF0IjoxNTYyMTkwNTI0LCJqdGkiOiJhYWFhYWFhYS1iYmJiLWNjY2MtZGRkZC1lZWVlZWVlZWVlZWUiLCJjbGllbnRfaWQiOiI1N2NiaXNoazRqMjRwYWJjMTIzNDU2Nzg5MCIsInVzZXJuYW1lIjoiamFuZWRvZUBleGFtcGxlLmNvbSJ9.CT-aGaJJCLRxDqEvMITntjF-RAbaVxDWgrfxPRYtgCD4l91uzNkfeuwbQFNbx-0eD7A86O8G-FwVqhBZXn1Wc12X-le_Ta_0Y4FX0YI3Zjqio1KSZIV2_AW5BUmd0FV4DMQplXWPMyJhJpYH7d8erDD8txAANMYL0EPA_pwDo1zQSlioas88ELF8sJa2dUMXLcOFnbqNkWQj_y_Y58iWRoVdwRz51pWzuT5NiJ4JOogolgBhYzwcWNwZE3VD7ht5y_udzDl4fHmadCxgNOacF4jVF91eLKH0PYZR6Dbf0a1Y7QDEYJQU3fYikoOVOL9DlszU5cRtitYbrrsaLrM_PQ"
)

func (f *PermissionsComponent) iHaveTheseRoles(rolesWriteJSON *godog.DocString) error {
	ctx := context.Background()
	roles := []models.Role{}
	m := f.MongoClient

	err := json.Unmarshal([]byte(rolesWriteJSON.Content), &roles)
	if err != nil {
		return err
	}

	for _, roleDoc := range roles {
		if err := f.putRolesInDatabase(ctx, m.Connection.Collection(m.ActualCollectionName(config.RolesCollection)), roleDoc); err != nil {
			return err
		}
	}

	return nil
}

func (f *PermissionsComponent) putRolesInDatabase(ctx context.Context, mongoCollection *dpMongoDriver.Collection, roleDoc models.Role) error {
	update := bson.M{
		"$set": roleDoc,
		"$setOnInsert": bson.M{
			"last_updated": time.Now(),
		},
	}
	_, err := mongoCollection.UpsertById(ctx, roleDoc.ID, update)
	if err != nil {
		return err
	}
	return nil
}

func (f *PermissionsComponent) iHaveThesePolicies(jsonInput *godog.DocString) error {
	ctx := context.Background()
	policies := []models.Policy{}
	m := f.MongoClient

	err := json.Unmarshal([]byte(jsonInput.Content), &policies)
	if err != nil {
		return err
	}

	for _, policy := range policies {
		if err := f.putPolicyInDatabase(ctx, m.Connection, policy, m.ActualCollectionName(config.PoliciesCollection)); err != nil {
			return err
		}
	}

	return nil
}

func (f *PermissionsComponent) putPolicyInDatabase(
	ctx context.Context,
	mongoConnection *dpMongoDriver.MongoConnection,
	policy models.Policy,
	collection string) error {
	update := bson.M{
		"$set": policy,
		"$setOnInsert": bson.M{
			"last_updated": time.Now(),
		},
	}
	_, err := mongoConnection.Collection(collection).UpsertById(ctx, policy.ID, update)
	if err != nil {
		return err
	}
	return nil
}

func (f *PermissionsComponent) adminJWTToken() error {
	err := f.APIFeature.ISetTheHeaderTo("Authorization", jwtAdminOnlyToken)
	return err
}

func (f *PermissionsComponent) publisherJWTToken() error {
	err := f.APIFeature.ISetTheHeaderTo("Authorization", jwtPublisherAndAdminToken)
	return err
}

func (f *PermissionsComponent) viewerJWTToken() error {
	err := f.APIFeature.ISetTheHeaderTo("Authorization", jwtViewerOnlyToken)
	return err
}

func (f *PermissionsComponent) basicUserJWTToken() error {
	err := f.APIFeature.ISetTheHeaderTo("Authorization", jwtBasicUserOnlyToken)
	return err
}

func (f *PermissionsComponent) publisherWithNoJWTToken() error {
	err := f.APIFeature.ISetTheHeaderTo("Authorization", "")
	return err
}
