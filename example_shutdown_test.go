// Copyright 2018 Kurt K, Tamás Gulácsi.
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package goracle_test

import (
	"database/sql"
	"log"

	"github.com/pkg/errors"

	goracle "gopkg.in/goracle.v2"
)

func ExampleShutdown() {
	dsn := "sys/sys as sysdba"
	db, err := sql.Open("goracle", dsn)
	if err != nil {
		log.Fatal(errors.Wrap(err, dsn))
	}
	defer db.Close()

	if err = exampleShutdown(db, goracle.ShutdownTransactionalLocal); err != nil {
		log.Fatal(err)
	}
}

func exampleShutdown(db *sql.DB, shutdownMode goracle.ShutdownMode) error {
	oraDB, err := goracle.DriverConn(db)
	if err != nil {
		return err
	}
	log.Printf("Beginning shutdown %v", shutdownMode)
	if err = oraDB.Shutdown(shutdownMode); err != nil {
		return err
	}
	// If we abort the shutdown process is over immediately.
	if shutdownMode == goracle.ShutdownAbort {
		return nil
	}

	log.Println("Closing database")
	if _, err = db.Exec("alter database close normal"); err != nil {
		return err
	}
	log.Println("Unmounting database")
	if _, err = db.Exec("alter database dismount"); err != nil {
		return err
	}
	log.Println("Finishing shutdown")
	if err = oraDB.Shutdown(goracle.ShutdownFinal); err != nil {
		return err
	}
	return nil
}
