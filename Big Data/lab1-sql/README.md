# DSGA1004 - BIG DATA
## Lab 1: SQLite
- Prof Brian McFee (bm106)
- Junge Zhang (jz3502)
- Jack Zhu (wz727)

					
*Handout date*: 2020-02-06

*Submission deadline*: 2020-02-19, 23:55 EST

### Goal of this exercise:
In this exercise, you will write queries to understand the concepts of relational databases and SQL, and write Python scripts to execute those queries.	
			

Database Provided:  `Sample_Song_Dataset.db`.  This can be obtained through the
`resources` section of the newclasses.nyu.edu website.

The database schema contains three tables:

- tracks
    - track_id
    - title
    - release
    - year
    - duration
    - artist_id
- artists
    - artist_id
    - artist_name
- artist_term
    - artist_id
    - term


Modify the given Python Script (`Lab1StarterCode.py`) to execute the following queries. Also, create a PDF file that includes for each question: its results (output) and an explanation of the query when appropriate.

In the assignment, you will need the sqlite3 package that can be installed by Ananconda using the command `conda install sqlite`.

To run your Script, please do not hard code your local database file/directory. Instead, make it the first argument for your command line (the code is included in the given Python Script). To run you Script, input `python Lab1StarterCode.py Sample_Song_Dataset` in your command line and execute it.  

	
1. Find id, name, and term of the artist who played the track with title "One Little Too Little".
2. Select all the unique track titles with the duration strictly shorter than `2` seconds. 
3. Find the ten longest (by duration) 10 track ids released between `2009` and `2014` (exclusive), ordered by decreasing duration.
4. Find the `20` least frequently used terms, ordered by increasing usage, and term in lexicographic order if terms are with the same frequency.
5. Find the artist name associated with the longest track duration.
6. Find the mean duration of all tracks.
7. Using only one query, find top `10` track ids, in lexicographic order,  whose artists strictly have more than `5` linked terms.
8. Index- Run Question 1 query in a loop for 100 times and note the minimum time taken. Now create an index on the column artist_id and compare the time. Share your findings in the report.
9. Find all tracks associated with artists that have the tag `eurovision winner` and delete them from the database, then roll back this query using a transaction.  Hint: you can select from the output of a select!

### Suggestions on submission
-After you complete each step, commit your changes to git so that we have a log of your progress.

-After printing an output of a question, also print out dashlines or leave spaces before you start next question. 
