# Lab 2 - Evaluation

---

In this assignment, you will build on your experience in the previous assignment, and conduct a more thorough evaluation of recommender systems using ranking metrics.

As always, please read through this entire assignment before getting started.

## Part 1: Datasets

The dataset that we'll be using is the "MovieLens-Latest-Small" set.  It is available in the shared folder on Prince under `/scratch/courses/DS-GA.3001-016-2020/ml-latest-small` and `/scratch/courses/DS-GA.3001-016-2020/ml-latest` .

In this assignment, you will use both the "small" and "full" versions of the MovieLens data, which are available in the shared folder on Prince .  As before, loading the data is done by the following:

```python
from lenskit.datasets import MovieLens

# Small data
data = MovieLens('/scratch/courses/DS-GA.3001-016-2020/ml-latest-small')

# Big data
data_full = MovieLens('/scratch/courses/DS-GA.3001-016-2020/ml-latest')
```

It is generally recommended that you complete the rest of the assignment first with the small data, and then repeat the process with the larger dataset.


## Part 2: Models

In addition to the two models you experimented with in lab 1, here you will also include a comparison to an implicit feedback model.

- Bias predictor: `lenskit.algorithms.basic.Bias`
- Biased matrix factorization: `lenskit.algorithms.als.BiasedMF`
- Implicit feedback matrix factorization: `lenskit.algorithms.als.ImplicitMF`

The first two models are designed to approximate explicit feedback (ratings).

The last model is designed to work with implicit feedback, so you should binarize the ratings data prior to fitting the model.
Use the rule that `R>=4` is a positive interaction and everything else is negative.
You should apply this rule prior to any data splitting.

Each model should be fit following the pattern from lab 1: use RNG seeds to ensure reproducibility, and a subset of training data should be held out for validation and hyper-parameter optimization.
Like before, use 20% of each user's interactions as test data via the `partition_users` method.

Unlike the previous assignment, we will now have multiple evaluation criteria available for hyper-parameter tuning.
The choice of which metric you use for model selection is up to you, but you must document your decisions in the report!

## Part 3: Evaluations

For each method, report the following performance metrics on the test set:

- MRR
- Prec@10
- Recall@10
- Prec@100
- Recall@100
- NDCG

As these metrics require binary relevance, you will need to binarize the ratings data here as well.  Again, use the `R>=4` rule to binarize ratings as positive.

Additionally, for the explicit feedback methods, report the RMSE and MAE metrics using the raw, explicit feedback of the test set.

## What to turn in

- The entire source code of your experiment (.py scripts and SLURM scripts)
- `REPORT.md`: a brief writeup of your experimental results as described above.

Your report should answer the following questions.

1. What selection criteria (choice of validation metric) did you use for model selection (hyper-parameter tuning) and why?
2. Would different model selection criteria have led to different outcomes?
3. Which method(s) perform best on the test for each of the metrics in question?

## General tips

- Complete the experiment first on the small dataset before proceeding to the large dataset.  You may want to write your scripts so that they accept the dataset name as a command line argument so that it is easy to reuse the same code for both.  You should expect to change your SLURM job request settings to allow for more memory and CPU time on the larger dataset.
- You may want to familiarize yourself with the batch and multi-eval utilities provided by lenskit.  It is not strictly necessary to use them, but they may simplify things considerably.
- In the large dataset, one train-test split will suffice.  You may use multiple splits on the smaller dataset if you wish.
- If you're feeling adventurous, you may want to investigate computing additional metrics which are not implemented in lenskit, such as AUC-ROC or average precision.  These metrics come with scikit-learn, but you would need to implement some wrapper functionality to use them here.
