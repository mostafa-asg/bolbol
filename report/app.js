const express = require("express");
const sqlite3 = require("sqlite3");
const util = require("util");

const app = express();
app.use(express.static("public"));
app.set("view engine", "ejs");
const port = 3000;
const dbFilename = "../report.db";

app.get("/report/:filename/words", (req, res) => {
  const filename = req.params.filename;
  var data = [];
  const db = new sqlite3.Database(dbFilename, sqlite3.OPEN_READONLY, err => {
    if (err) {
      console.log("error: cannot open database");
    } else {
      var query = `
      SELECT DISTINCT word,
        (select count(*) from word_reports where filename = '%s' AND answered = 1 and word=T1.word) as success,
        (select count(*) from word_reports where filename = '%s' AND answered = 0 and word=T1.word) as failure
      FROM word_reports AS T1 where filename = '%s' ORDER BY failure DESC, SUCCESS ASC;
      `;

      query = util.format(query, filename, filename, filename);

      db.all(query, (err, rows) => {
        if (err) {
          console.log("error: could not execute query", err);
        } else {
          rows.forEach(row => {
            data.push({
              word: row.word,
              filename: row.filename,
              success: row.success,
              failure: row.failure
            });
          });

          db.close();
          const model = {
            data: data
          };

          res.render("wordsReport", model);
        }
      });
    }
  });
});

app.get("/", (req, res) => {
  var data = [];
  const db = new sqlite3.Database(dbFilename, sqlite3.OPEN_READONLY, err => {
    if (err) {
      console.log("error: cannot open database");
    } else {
      const query = `
            select distinct filename, 
            (select count(*) from reports as TW where 
                                  TW.filename=T1.filename AND 
                                  T1.date <= date('now') AND date >= date('now','-7 day')) as weekCount, 
            (select count(*) from reports as TW where 
                                  TW.filename=T1.filename AND 
                                  T1.date <= date('now') AND date >= date('now','-30 day')) as monthCount,
            (select count(*) from reports as T0 where 
                                  T0.filename=T1.filename) as totalCount
            from reports as T1;
      `;

      db.all(query, (err, rows) => {
        if (err) {
          console.log("error: could not execute query", err);
        } else {
          rows.forEach(row => {
            data.push({
              filename: row.filename,
              weekCount: row.weekCount,
              monthCount: row.monthCount,
              totalCount: row.totalCount
            });
          });

          db.close();
          const model = {
            data: data
          };

          res.render("index", model);
        }
      });
    }
  });
});

app.listen(port, () => console.log(`listening on port ${port}`));
