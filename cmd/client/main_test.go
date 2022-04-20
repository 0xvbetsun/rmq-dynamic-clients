// Entry point for application's clients
package main

import "testing"

func Test_parseCmd(t *testing.T) {
	tests := []struct {
		name    string
		cmd     string
		wantCmd string
		wantArg string
		wantErr bool
	}{
		{name: "command with 1 arg", cmd: "AddItem abc", wantCmd: "AddItem", wantArg: "abc"},
		{name: "command with no args", cmd: "GetAllItems", wantCmd: "GetAllItems"},
		{name: "without command", cmd: "", wantErr: true},
		{name: "command with more than expected args", cmd: "AddItem abc bcd cde", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotArg, err := parseCmd(tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCmd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCmd != tt.wantCmd {
				t.Errorf("parseCmd() gotCmd = %v, want %v", gotCmd, tt.wantCmd)
			}
			if gotArg != tt.wantArg {
				t.Errorf("parseCmd() gotArg = %v, want %v", gotArg, tt.wantArg)
			}
		})
	}
}
