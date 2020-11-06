# Lab 1 - Building a movie recommender

---

In this assignment, you will benchmark some of the methods that we've discussed in class on a small dataset from [MovieLens](https://grouplens.org/datasets/movielens/).

You'll be using the implementations provided by the [LensKit](https://lkpy.readthedocs.io/en/stable/index.html) package, which should already be installed in your conda environment from lab 0.  If you did not complete lab 0, now would be a good time to go back and do it!

Please read through this entire assignment before getting started.

## Part 1: Loading the data

The dataset that we'll be using is the "MovieLens-Latest-Small" set.  IT is available in the shared folder on Prince under `/scratch/courses/DS-GA.3001-016-2020/ml-latest-small`.

LensKit provides data loader classes to interface with different versions of the MovieLens data: older versions use a structured `.dat` format, while the recent versions (like the one we're using) come in .csv format.  Loading the data is done by the following:

```python
from lenskit.datasets import MovieLens

data = MovieLens('/scratch/courses/DS-GA.3001-016-2020/ml-latest-small')
```

After executing the above code, you should be able to access the ratings data as a pandas DataFrame:

```python
>>> data.ratings
        user    item  rating   timestamp
0          1       1     4.0   964982703
1          1       3     4.0   964981247
2          1       6     4.0   964982224
3          1      47     5.0   964983815
4          1      50     5.0   964982931
...      ...     ...     ...         ...
100831   610  166534     4.0  1493848402
100832   610  168248     5.0  1493850091
100833   610  168250     5.0  1494273047
100834   610  168252     5.0  1493846352
100835   610  170875     3.0  1493846415

[100836 rows x 4 columns]
```

## Part 2: Splitting the data

To benchmark recommendation algorithms, you will need to partition (split) the data to simulate the effect of having unobserved interactions.

Data partitioning for recommender systems requires a bit more care than in the typical independent-and-identically-distributed (IID) setting you may be familiar with in machine learning.
The quantity that we want to evaluate is the accuracy of the recommender *on average across users*.
Computing this requires that we have a set of *test users* who have some observed interaction history to learn from, and some interactions held out.

LensKit provides user-conditional data partitioning via the following pattern:

```python
from lenskit.crossfold import partition_users, SampleFrac

# How many partitions we'll evaluate on
N_SPLITS = 5

# Portion of each test user's history to hold out
FRAC_SPLIT = 0.2

# Carve out just the necessary data from the dataframe
ratings = data.ratings[['user', 'item', 'rating']]

for train, test in partition_users(ratings, N_SPLITS, SampleFrac(FRAC_SPLIT)):
    # fit the model(s) on train
    # evaluate the model(s) on test
```
The above example would provide `N_SPLITS=5` different partitions of the data into `train` and `test` sets.

Because this dataset is relatively small, your evaluation should consist of at least 5 trials (indepdendent partitions).

## Part 3: Fitting and tuning models

Compare the following methods:

- Bias predictor: `lenskit.algorithms.basic.Bias` (item bias only)
- Biased matrix factorization: `lenskit.algorithms.als.BiasedMF`

For each model, there are hyper-parameters that you can tune to improve (hopefully!) the performance of the model.
For the bias predictor, you may want to experiment with different settings of the damping factor.
For matrix factorization, the regularization and dimensionality (`features` in the API) will be important.

To tune the parameters of your models, it is recommended to further partition the training set into `train` and `validation`.
Choose the hyper-parameter setting that performs best (by RMSE) on the validation set, and then record its performance on the test set.
This must be done independently for each trial.

## Part 4: Comparing models

For each configuration of each model, report the "root mean squared error" (RMSE) of its predictions for your train, validation, and test splits.  This should be done by calling the `predict` method of the fit model on the list of unobserved `(user, item)` interactions contained in the test set, and evaluating the estimated ratings against the truth:

```python
from lenskit.metrics.predict import rmse

# Fit the model to the training set
model.fit(train)

# Make predictions on held-out interactions
ratings_est = model.predict(test[['user', 'item']])

# Compare predictions to the held-out ratings
score = rmse(ratings_est, test['rating'])
```

Report these quantities independently for each trial, as well as an aggregate summary (mean and standard deviation) across trials.


## What to turn in

- The entire source code of your experiment (.py scripts and SLURM scripts)
- `REPORT.md`: a brief writeup of your experimental results as described above.

## General tips

- When partitioning data, provide a seed by saying `rng_spec=` in the call to `partition_users` and `BiasedMF`.  This will ensure that your data partitioning and model fitting is reproducible across runs, and help with debugging.  **Always seed your RNGs!**
- This dataset is small enough that you should be able to develop your code locally on a laptop, at least up until the parameter optimization.  It's fine to do so (recommended even!), but make sure that your final submission runs as a submitted job on the HPC.
- Get the entire system working end-to-end first with the `Bias` model, and then go back to add the `BiasedMF` model.
- Make sure that each trial uses a fresh instantiation of the recommendation algorithm, otherwise you may observe skewed results!
- Refer back to the "lab 0" repository for environment configuration and reference material for submitting jobs to HPC.
- You may want to precompute your data splits and save them to temporary files that are loaded later on.  This is fine, just make sure that your RNG is seeded so that the results are reproducible.
- The range of hyper-parameter settings you use for each model is up to your discretion, but don't go overboard with it!  The important thing here is to get a sense of the range of feasible settings, not to precisely optimize the model(s) for this particular dataset.
