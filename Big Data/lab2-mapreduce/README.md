# DSGA1004 - BIG DATA
## Lab 2: Hadoop
- Prof Brian McFee (bm106)
- Weicheng Zhu (wz727)
- Junge Zhang (jz3502)

					
*Handout date*: 2020-02-20

*Submission deadline*: 2020-03-04, 23:55 EST

## 0. Requirements

Sections 1 and 2 of this assignment are designed to get you familiar with HPC and the work-flow of running Hadoop jobs.

For full credit on this lab assignment, you will need to provide working solutions for sections 4 (matrix multiplication) and 5 (table joins).
You are provided with small example inputs for testing these programs, but we will run your programs on larger data for grading purposes.
Be sure to commit your final working implementations to git (and push your changes) before the submission deadline!


## 1. High-performance Computing (HPC) at NYU

This lab assignment will require the use of the
Hadoop cluster run by the NYU high-performance
computing (HPC) center.  To learn more about HPC at
NYU, refer to the [HPC Wiki](https://wikis.nyu.edu/display/NYUHPC/High+Performance+Computing+at+NYU).

By now, you should have received notification at
your NYU email address that your HPC account is active. If you have not received this notification yet, please contact the instructor immediately.

If you're new to HPC, please read through the
[tutorial](https://wikis.nyu.edu/display/NYUHPC/Tutorials) section of the wiki, and in particular, the [Getting started on Dumbo](https://wikis.nyu.edu/display/NYUHPC/Clusters+-+Dumbo) section.

Logging in Dumbo on Linux or Mac from the NYU network:
```bash
ssh netid@dumbo.hpc.nyu.edu
```
Uploading file to Dumbo:
```bash
scp local_dir netid@dumbo.hpc.nyu.edu:dumbo_dir
```
Downloading file from Dumbo:
```bash
scp netid@dumbo.hpc.nyu.edu:dumbo_dir local_dir
```

While it is possible to transfer files directly to and from Dumbo, we strongly recommend that you use git (and GitHub) to synchronize your code rather than `scp`.  This way, you can be sure that your submitted project is always up to date with the code being run on the HPC.  To do this, you may need to set up a new SSH key (on Dumbo) and add it to your GitHub account; instructions for this can be found [here](https://help.github.com/en/enterprise/2.15/user/articles/adding-a-new-ssh-key-to-your-github-account).

**Note**: Logging into the HPC from outside the NYU
network can be somewhat complicated.  Instructions
are given
[here](https://wikis.nyu.edu/display/NYUHPC/Logging+in+to+the+NYU+HPC+Clusters).



## 2. Hadoop and Hadoop-streaming

In lecture, we discussed the Map-Reduce paradigm in the abstract, and did not dive into the details of the Hadoop implementation.  Hadoop is an open-source implementation of map-reduce written in Java.
In this lab, you will be implementing map-reduce jobs using Hadoop's "streaming" mode, which allows mappers and reducers to be implemented in any programming language.  In our case, we'll be using Python.


### Environment setup

To setup required environment, you need to execute the following command under the git repo directory every time logging in Dumbo:
```bash
source shell_setup.sh
```
These modifications add shortcuts for interacting with the Hadoop distributed filesystem (`hfs`) and launching map-reduce jobs (`hjs`).

*Note*: For convinence, feel free to copy-paste the contents of that file into `~/.bashrc` so that you don't need to re-run setup everytime you log in.

### Github on Dumbo

Follow the [instruction](https://help.github.com/en/github/authenticating-to-github/connecting-to-github-with-ssh) to clone Github repo through SSH (cloning via https link may encounter problems with HPC). The repository includes three problems on MapReduce:

## 3. A first map-reduce project (Not for grading)

Included within the repository under `word_count/` is a full implementation of the "word-counting" program, and an example input text file (`book.txt`).

The program consists of four files:
```
src/mapper.py
src/mapper.sh
src/reducer.py
src/reducer.sh
```
Why four files?  Clearly, two of them relate to the *mapper* and two relate to the *reducer*.
The reason for this is that we use the shell scripts (`.sh` extensions) to load dependencies on the worker nodes of the cluster before the mapper or reducer functions (`.py` extensions) are executed.
In this case, the only dependency is `Python`.
Once the dependency modules are loaded, the shell script (`mapper.sh` or `reducer.sh`) executes the mapper or reducer, which read from the standard input and write to the standard output.
Hadoop-streaming will coordinate the communication between mappers and reducers for us.


### Testing the mapper and reducer with Python using CPUs

Before we move on, it's a good idea to run these programs locally so we know what to expect.  (*Hint*: this is also an easy way to debug, as long as you have a small input on hand!)

You can run the mapper by going into the `word_count/src/` directory, and running the following commands:
```bash
cat ../book.txt | python3 mapper.py
```
The first command enables Python in your environment so that it matches what will run on the cluster.
The second command will run the contents of `book.txt` through `mapper.py`, resulting in an unordered list of `(key, value)` pairs.

These key-value pairs will be in the same order as they appear in `book.txt`, but Map-Reduce will want to sort them by key to make grouping easier.
To simulate this, run the same command followed by `sort`:
```bash
cat ../book.txt | python3 mapper.py | sort
```
This will produce the same output as above, but now you should see all repetitions of the same word grouped together.

Finally, you can run this through the `reducer` by adding one more step:
```bash
cat ../book.txt | python3 mapper.py | sort | python3 reducer.py
```
After running this command, you should see the total counts of each word in `book.txt`!
Remember, we did this all on one machine without using Hadoop, but you should now have a sense of what the inputs and outputs of these programs look like.


### Launching word-count on Hadoop cluster

Before we can launch the word counter on the cluster, we will need to place the input file in the Hadoop distributed file system (HDFS).

To do this, issue the following command from inside the `word_count` directory:

* Put text file on HDFS of Hadoop cluster
```bash
hfs -put book.txt
```
* Launch map-reduce job for counting words on Hadoop cluster 
This is done by issuing the following command from inside the `word_count` directory:
```bash
hjs -file src/ -mapper src/mapper.sh -reducer src/reducer.sh -input book.txt -output word_count.out
```

The `hjs` command submits the job; the `-file` parameter indicates which files will be distributed to the worker nodes (i.e., the source code); the `-mapper` and `-reducer` parameters specify the paths to the mapper and reducer scripts; and the `-input` and `-output` paths specify the input file(s) and output file paths for the job.

* After the job finishes, check HDFS (the output of Hadoop jobs will be stored in HDFS)
```bash
hfs -ls
```
and you should see two files: `book.txt` and `word_count.out`, the latter being a directory.
If you run `hfs -ls word_count.out/` you will see several file "parts", each corresponding to a single reducer node.
To retrieve the results of the computation, run
```bash
hfs -get word_count.out 
```
to get all the partial outputs.
To get the complete output as one file, run:
```bash
hfs -getmerge word_count.out word_count_total.out
```
After running these commands, the results of the computation will be available to you through the usual unix file system.

At this point, you should now have successfully run your first Hadoop job!


## 4. Matrix products
In the next part of this assignment, you will develop a map-reduce algorithm for matrix multiplication.  Recall that the product of two matrices `A * B` where `A` has shape `(n, k)` and `B` has shape `(k, m)` is a matrix `C` of shape `(n, m)` where for each `i` and `j` coordinates,
```
C[i, j] = A[i, 1] * B[1, j] + A[i, 2] * B[2, j] + ... + A[i, k] * B[k, j]
```

In the `matmul/` directory, you will find skeleton code for implementing matrix multiplication `matmul/src/`.
Accompanying this is an example pair of matrices to multiply: `matmul/small/A.txt` and `matmul/small/B.txt`.
In this assignment, you can assume that `A` and `B` will always be of compatible dimension, and that `A` will be on the left side of the multiplication.
The mapper program should take as command-line arguments the dimensions of the product `C`, that is, the number of rows of `A` and the number of columns of `B`.


Your job in this part of the assignment is to fill in the skeleton code (`mapper.py` and `reducer.py`) by designing the intermediate key-value format, and implementing the rest of the program accordingly.

Note that the data now comes spread into two files (`A.txt` and `B.txt`), which must both be processed by the mapper. You shall have the output like `output.txt` as the product of matrix `A` and `B`. 

**Note**: the input and output files should be in the sparse format of the matrix (i.e. only nonzero entries of matrix `M` are recorded as tuples `i,j,M[i, j]`)

e.g. Any example of input and output are following: 
```
A:  0,0,1.0  B:  0,1,2.0    A*B:   0,0,6.0
    0,1,2.0      1,0,3.0           0,1,2.0
    1,0,3.0                        1,1,6.0
```

In Hadoop streaming, each file can be spread across multiple mappers, and if a directory is given as the `-input` option, then all files are processed.
Within the mapper, the identity of the file currently being processed can be found in the environment variable `mapreduce_map_input_file`, as noted in the code comments.  (Hint: this will be useful in part 5 below!)




To run the job, place the data on HDFS by issuing the command:
```bash
hfs -put small
```
from within the `matmul` directory.
Once you have implemented your mapper and reducer, the job should be launched by the following command:
```bash
hjs -files src/ -mapper src/mapper.sh -reducer src/reducer.sh -input small/ -output small_product
```
from within the `matmul` directory.
The output of the computation can be retrieved using `hfs -getmerge` like in the example section above.


**Note**: because this program will use multiple input files and environment variables to track which file is being processed, it is a bit more difficult to run locally.  You can test it as follows:
```bash
mapreduce_map_input_file=A.txt python3 src/mapper.py 2 3 < small/A.txt > tmp_A.txt
mapreduce_map_input_file=B.txt python3 src/mapper.py 2 3 < small/B.txt > tmp_B.txt
sort tmp_A.txt tmp_B.txt | python3 src/reducer.py
```
which will temporarily set the input file environment variable while running the mapper on each input (lines 1 and 2), and run the reducer on the sorted intermediate values.

### Why the "2 3" parameters?

The mapper needs to know the shape of the matrix that you will eventually produce when multiplying `A * B`: specifically, the number of rows in `A` (2-by-4 for the small example) and the number of columns in `B` (4-by-3).
You will need to specify these values as parameters to the mapper script as shown above.

## 5. Tabular data

In the final section of the lab, you are given two data files in comma-separated value (CSV) format.
These data files (`joins/music_small/artist_term.csv` and `joins/music_small/track.csv`) contain the same music data from the previous lab assignment on SQL and relational databases.  Specifically, the file `artist_term.csv` contains data of the form
```
ARTIST_ID,tag string
```
and `track.csv` contains data of the form
```
TRACK_ID,title string,album string,year,duration,ARTIST_ID
```

No skeleton code is provided for this part, but feel free to adapt any code from the previous sections that you've already completed.

### 5.1 Joining tables

For the first part, implement a map-reduce program which is equivalent to the following SQL query:
```
SELECT 	track.artist_id, track.track_id, artist_term.tag
FROM	track INNER JOIN artist_term 
ON 	track.artist_id = artist_term.artist_id
WHERE   track.year > 1990
```

The program should be executable in a way similar to the matrix multiplication example, for example:
```bash
hjs -files src/ -mapper src/join_mapper.sh -reducer src/join_reducer.sh -input music_small/ -output join_query
```

### 5.2 Aggregation queries

For the last part, implement a map-reduce program which is equivalent to the following SQL query:
```
SELECT 	artist_term.tag, min(track.year), avg(track.duration), count(artist_term.artist_id)
FROM	track LEFT JOIN artist_term
ON	track.artist_id = artist_term.artist_id
GROUP BY artist_term.tag
```
That is, for each artist ID, compute the maximum year of release, average track duration and the total number of terms matching the artist.  **Note**: the number of terms for an artist could be zero!

The program should be executable by the following command:
```bash
hjs -files src/ -mapper src/group_mapper.sh -reducer src/group_reducer.sh -input music_small/ -output group_query
```
