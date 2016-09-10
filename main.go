package main

import (
	"github.com/kataras/iris"
	"github.com/tarantool/go-tarantool"
	"log"
	"time"
)


func main() {
	server := "77.244.214.4:3301"
	spaceNo := "sessions"
	indexNo := uint32(0)
	opts := tarantool.Opts{
		Timeout:       150 * time.Millisecond,
		Reconnect:     1 * time.Second,
		MaxReconnects: 3,
		User:          "goClient",
		Pass:          "TeddyTarantoolS1cret",
	}

	client, err := tarantool.Connect(server, opts)
	if err != nil {
		log.Fatalf("Failed to connect: %s", err.Error())
	}
	resp, err := client.Ping()
	log.Println(resp.Code)
	log.Println(resp.Data)
	log.Println(err)
	// insert new tuple { 10, 1 }
	resp, err = client.Insert(spaceNo, []interface{}{uint(1), 1})
	// or
	resp, err = client.Insert("sessions", []interface{}{uint(2), 1})
	log.Println("\n\nInsert")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// delete tuple with primary key { 10 }
	resp, err = client.Delete(spaceNo, indexNo, []interface{}{uint(1)})
	// or
	resp, err = client.Delete("sessions", "primary", []interface{}{uint(10)})
	log.Println("\n\nDelete")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// replace tuple with { 13, 1 }
	resp, err = client.Replace(spaceNo, []interface{}{uint(2), 1})
	// or
	resp, err = client.Replace("sessions", []interface{}{uint(2), 1})
	log.Println("\n\nReplace")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// update tuple with primary key { 13 }, incrementing second field by 3
	resp, err = client.Update(spaceNo, indexNo, []interface{}{uint(13)}, []interface{}{[]interface{}{"+", 1, 3}})
	// or
	resp, err = client.Update("sessions", "primary", []interface{}{uint(13)}, []interface{}{[]interface{}{"+", 1, 3}})
	log.Println("\n\nUpdate")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// insert tuple {15, 1} or increment second field by 1
	resp, err = client.Upsert(spaceNo, []interface{}{uint(15), 1}, []interface{}{[]interface{}{"+", 1, 1}})
	// or
	resp, err = client.Upsert("sessions", []interface{}{uint(15), 1}, []interface{}{[]interface{}{"+", 1, 1}})
	log.Println("\n\nUpsert")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// select just one tuple with primay key { 15 }
	resp, err = client.Select(spaceNo, indexNo, 0, 1, tarantool.IterEq, []interface{}{uint(15)})
	// or
	resp, err = client.Select("sessions", "primary", 0, 1, tarantool.IterEq, []interface{}{uint(15)})
	log.Println("\n\nSelect")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// select tuples by condition ( primay key > 15 ) with offset 7 limit 5
	// BTREE index supposed
	resp, err = client.Select(spaceNo, indexNo, 7, 5, tarantool.IterGt, []interface{}{uint(15)})
	// or
	resp, err = client.Select("sessions", "primary", 7, 5, tarantool.IterGt, []interface{}{uint(15)})
	log.Println("\n\nSelect")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// call function 'func_name' with arguments
	resp, err = client.Call("func_name", []interface{}{1, 2, 3})
	log.Println("\n\nCall")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)

	// run raw lua code
	resp, err = client.Eval("return 1 + 2", []interface{}{})
	log.Println("\n\nEval")
	log.Println("Error", err)
	log.Println("Code", resp.Code)
	log.Println("Data", resp.Data)
	api := iris.New()
	api.Get("/hi", hi)
	api.Listen(":8080")
}

func hi(ctx *iris.Context){
	ctx.Write("Hi %s", "iris")
}
