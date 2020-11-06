#!/usr/bin/env python
# -*- encoding: utf-8 -*-

# USAGE:
#   python Lab1.py Sample_Song_Dataset.db

import sys
import sqlite3


# The database file should be given as the first argument on the command line
# Please do not hard code the database file!
db_file = sys.argv[1]
    

# We connect to the database using 
with sqlite3.connect(db_file) as conn:
    # We use a "cursor" to mark our place in the database.
    # We could use multiple cursors to keep track of multiple
    # queries simultaneously.
    cursor = conn.cursor()

    # This query counts the number of tracks from the year 1998
    year = ('1998',)
    cursor.execute('SELECT count(*) FROM tracks WHERE year=?', year)

    # Since there is no grouping here, the aggregation is over all rows
    # and there will only be one output row from the query, which we can
    # print as follows:
    #print('Tracks from {}: {}'.format(year[0], cursor.fetchone()[0]))
    # The [0] bits here tell us to pull the first column out of the 'year' tuple
    # and query results, respectively.

    # ADD YOUR CODE STARTING HERE
   # cursor.execute('DROP INDEX artist_index_1')
    #cursor.execute('DROP INDEX artist_index_2')
    #cursor.execute('DROP INDEX artist_index_3')
    print("---------- problem 1 ----------")
    cursor.execute('SELECT a.artist_id, a.artist_name, ar.term FROM artists a, artist_term ar, tracks t where t.artist_id = a.artist_id and ar.artist_id = a.artist_id and t.title = "One Little Too Little";')
    print(cursor.fetchall())
    print("---------- problem 2 ----------")
    cursor.execute('SELECT distinct t.title FROM tracks t where t.duration<2;')
    print([ o[0] for o in cursor.fetchall()])
    print("---------- problem 3 ----------")
    cursor.execute('SELECT t.track_id FROM tracks t where t.year>2009 and t.year<2014 order by t.duration desc limit 10;')
    print([o[0] for o in cursor.fetchall()])
    print("---------- problem 4 ----------")
    cursor.execute('SELECT term from artist_term group by term order by count(artist_id),term limit 20;')
    print([o[0] for o in cursor.fetchall()])
    print("---------- problem 5 ----------")
    cursor.execute('SELECT a.artist_name from artists a,tracks t where a.artist_id = t.artist_id order by t.duration desc limit 1;')
    print(cursor.fetchall()[0][0])
    print("---------- problem 6 ----------")
    cursor.execute('select avg(duration) from (SELECT duration from tracks group by track_id);')
    print(cursor.fetchall()[0][0])
    print("---------- problem 7 ----------")
    cursor.execute('select t.track_id from tracks t,(select a.artist_id from artist_term a group by a.artist_id having count(distinct a.term) > 5) temp where temp.artist_id = t.artist_id order by t.track_id limit 10;')
    print([o[0] for o in cursor.fetchall()])
    #xxx dont comment out xxx Below generates wrong result but don't know the reason
    #cursor.execute('select t.track_id, from artist_term a inner join tracks t on a.artist_id = t.artist_id group by a.artist_id having count(distinct a.term) > 5 order by t.track_id limit 10;')
    #print(cursor.fetchall())
    print("---------- problem 8 ----------")
    import time
    duration1 = []
    for i in range(100):
        s = time.time()
        cursor.execute('SELECT a.artist_id, a.artist_name, ar.term FROM artists a, artist_term ar, tracks t where t.artist_id = a.artist_id and ar.artist_id = a.artist_id and t.title = "One Little Too Little";')
        d = time.time()-s
        duration1.append(d)
     #   if(i%20 == 0):
      #      print("----- time to execute #1 query in iteration ",i,"is",d,"-----")
    print("minimum time to run #1 query before creating index on the column artist_id:",min(duration1))
    cursor.execute('CREATE INDEX artist_index_1 ON artist_term (artist_id);')
    cursor.execute('CREATE INDEX artist_index_2 ON artists (artist_id);')
    cursor.execute('CREATE INDEX artist_index_3 ON tracks (artist_id);')
    duration2 = []
    code = '''
def input(): 
    cursor.execute('SELECT a.artist_id, a.artist_name, ar.term FROM artists a, artist_term ar, tracks t where t.artist_id = a.artist_id and ar.artist_id = a.artist_id and t.title = "One Little Too Little";')   
'''
    import timeit
    for i in range(100):
        s = timeit.timeit(stmt = code,number = 1)
        duration2.append(s)
     #   if(i%20 == 0):
      #      print("----- time to execute #1 query in iteration ",i,"is",s,"-----")
    print("minimum time to run #1 query after creating index on the column artist_id:",min(duration2))
    print("---------- problem 9 ----------")
    cursor.execute('select count(DISTINCT track_id) from tracks;')
    print("number of unique track id before the delete and rollback:",cursor.fetchall()[0])
    cursor.execute('select distinct track_id from tracks t inner join (SELECT distinct artist_id from artist_term where term = "eurovision winner") temp on temp.artist_id=t.artist_id;')
    print("The tracks associated with artists that have the tag eurovision winner")
    print([o[0] for o in cursor.fetchall()])
    cursor.execute('BEGIN TRANSACTION problem;')
    cursor.execute('delete from tracks where track_id in (select track_id from tracks t inner join (SELECT distinct artist_id from artist_term where term = "eurovision winner") temp on temp.artist_id=t.artist_id);')
    cursor.execute('select count(DISTINCT track_id) from tracks;')
    print("number of unique track id after the delete before the rollback:",cursor.fetchall()[0])
    cursor.execute('ROLLBACK TRANSACTION problem;')
    cursor.execute('select count(DISTINCT track_id) from tracks;')
    print("number of unique track id after the delete and rollback",cursor.fetchall()[0])