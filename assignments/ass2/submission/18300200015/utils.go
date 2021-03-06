package main
import (
	"fmt"
	"time"
	
	// YOUR CODE BEGIN remove the follow packages if you don't need them
	"reflect"
	"sync"
	// YOUR CODE END

	_ "github.com/go-sql-driver/mysql"
	sql "github.com/jmoiron/sqlx"
)

var (
	// YOUR CODE BELOW
	EvaluatorID   = "18300200015" // your student id, e.g. 18307130177
	SubmissionDir = "../../../ass1/submission/" // the relative path the the submission directory of assignment 1, it should be "../../../ass1/submission/"
	User          = "root" // the user name to connect the database, e.g. root
	Password      = "v8k2s6d5n1c3g@"//omit // the password for the user name, e.g. xxx
	// YOUR CODE END
)

// ConcurrentCompareAndInsert is similar with compareAndInsert in `main.go`, but it is concurrent and faster!
func ConcurrentCompareAndInsert(subs map[string]*Submission) {
	start := time.Now()
	defer func() {
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/ass1_result_evaluated_by_%s", User, Password, EvaluatorID))
		if err != nil {
			panic(nil)
		}
		rows, err := db.Query("SELECT COUNT(*) FROM comparison_result")
		if err != nil {
			panic(err)
		}
		rows.Next()
		var cnt int
		err = rows.Scan(&cnt)
		if err != nil {
			panic(err)
		}
		if cnt == 0 {
			panic("ConcurrentCompareAndInsert Not Implemented")
		}
		fmt.Println("ConcurrentCompareAndInsert takes ", time.Since(start))
	}()
	// YOUR CODE BEGIN
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/ass1_result_evaluated_by_%s", User, Password, EvaluatorID))
	if err != nil {
			fmt.Println("error1")
			panic(nil)
		}
	wg := sync.WaitGroup{}
	for submitter, sub := range subs {
		for comparer, sub2 := range subs {
			for i := 0; i < NumSQL; i++ {
				wg.Add(1)
				go func(i int,submitter string,comparer string,sub *Submission,sub2 *Submission){
				var equal int
				if reflect.DeepEqual(sub.sqlResults[i], sub2.sqlResults[i]) {
					equal = 1
				} else {
					equal = 0
				}
				s := fmt.Sprintf("INSERT INTO comparison_result VALUES ('%s', '%s', %d, %d)", submitter, comparer, i+1, equal)
				_, err := db.Exec(s)
				if err != nil {
					fmt.Println(s)
					panic(err)
				}
				wg.Done()
				}(i,submitter,comparer,sub,sub2)
			}
			wg.Wait()	
		}
}
	
	
	// YOUR CODE END
}

// GetScoreSQL returns a string which contains only ONE SQL to be executed, which collects the data in table
// `comparision_result` and inserts the score of each submitter on each query into table `score`
func GetScoreSQL() string {
	var SQL string
	SQL = "SELECT 1" // ignore this line, it just makes the returned SQL a valid SQL if you haven't written yours.
	// YOUR CODE BEGIN
		SQL=`INSERT INTO score(submitter,item,score,vote) 
(SELECT result.submitter,result.item,1 AS score,vote
	FROM (SELECT submitter,item,count(*)AS vote
	FROM comparison_result
	WHERE is_equal = 1
	GROUP BY submitter, item) AS result
    WHERE vote >= ALL(SELECT count(*) FROM comparison_result AS other 
					WHERE is_equal=1 AND result.item=other.item GROUP BY other.submitter,other.item))
UNION
(SELECT result2.submitter,result2.item,0 AS score,vote
	FROM (SELECT submitter,item,count(*)AS vote
	FROM comparison_result
	WHERE is_equal = 1
	GROUP BY submitter, item) AS result2
    WHERE vote < some(SELECT count(*) FROM comparison_result AS other2 
					WHERE is_equal=1 AND result2.item=other2.item GROUP BY other2.submitter,other2.item));`

	// YOUR CODE END
	return SQL
}

func GetScore(db *sql.DB, subs map[string]*Submission) {
	// YOUR CODE BEGIN
		rows,err:=db.Query("SELECT submitter,item,score from score")
	if err != nil {
        fmt.Println(err)
	}
	var id string
	var item,score int
	for rows.Next(){
		rows.Scan(&id,&item,&score)
		subs[id].score[item]=score
	}


	// YOUR CODE END
}

