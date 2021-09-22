package client

import (
	"testing"

	"github.com/sirupsen/logrus"
)

var client = &Client{
	host: "http://api.pts.se/PTSNumberService/Pts_Number_Service.svc/json/SearchByNumber",
	log:  logrus.New(),
}

func TestClient_GetOperatorDetails(t *testing.T) {
	tests := []struct {
		name    string
		msisdn  string
		want    string
		wantErr bool
	}{
		{
			name:   "valid swedish phone number(+46) should return correct operator name",
			msisdn: "+46107500500",
			want:   "Telness AB",
		},
		{
			name:   "valid swedish phone number(without +46) should return correct operator name",
			msisdn: "0107500500",
			want:   "Telness AB",
		},
		{
			name:   "invalid swedish phone number should return missing operator",
			msisdn: "01050",
			want:   "Operatör saknas",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetOperatorDetails(tt.msisdn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetOperatorDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.D.Name != tt.want {
				t.Errorf("Client.GetOperatorDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}
