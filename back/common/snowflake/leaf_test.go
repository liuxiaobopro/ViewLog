package snowflake

import (
	"fmt"
	"testing"
)

func TestNewLeaf(t *testing.T) {
	type args struct {
		dataCenterId int64
		workerId     int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				dataCenterId: 1,
				workerId:     1,
			},
		},
		{
			name: "test1",
			args: args{
				dataCenterId: 1,
				workerId:     2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leaf, err := NewLeaf(tt.args.dataCenterId, tt.args.workerId)
			if err != nil {
				panic(err)
			}

			for i := 0; i <= 1; i++ {
				id, err := leaf.NextId()
				if err != nil {
					panic(err)
				}
				fmt.Printf("id = %d \n", id)
				// fmt.Printf("64bit id = %064b \n", id)
				// fmt.Printf("string id = %s \n", strconv.FormatInt(id, 10))
			}
		})
	}
}
