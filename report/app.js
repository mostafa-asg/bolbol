const express = require("express");
var sqlite3 = require("sqlite3");

const app = express();
app.use(express.static("public"));
app.set("view engine", "ejs");
const port = 3000;

app.get("/", (req, res) => {
  var data = [];
  const db = new sqlite3.Database(
    "../report.db",
    sqlite3.OPEN_READONLY,
    err => {
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
    }
  );
});

app.listen(port, () => console.log(`listening on port ${port}`));
