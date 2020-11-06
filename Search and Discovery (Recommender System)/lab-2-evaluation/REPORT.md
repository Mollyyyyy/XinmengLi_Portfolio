# Lab 2 - Report

Name: Xinmeng Li
NetID: xl1575

## Description of experiment
### Steps
1. Create a binarized rating colume in data, where 1 for rating>=4 and otherwise 0.
2. First split the data into train_val and test with FRAC_SPLIT = 0.2, rng seed = 0, then split train_val into train and val with FRAC_SPLIT = 0.2, rng seed = 0. 
3. Use MultiEval to fit models on train, tune hyperparameter(damping 10,100,200 for bias, and features 20,100,200 for biasedmf and implicitmf), and select the optimal hyperparameter where the ndcg on validation is the highest. Here we use original rating for explicit feedback model and binary rating for implicit feedback model when fitting the model. 
4. Fit the model with optimal setting on train_val and measure the metrics on test set.\
Note that for explicit feedback models, the score in recommendations.parquet is binarized before the binary-relevance metric evaluation; for implicit feedback models, the original score in recommendations.parquet is preserved, since the implicitmf fits on the binary train_val dataset.
## Results
### Small dataset
AlgoClass damping   nrecs  recip_rank      ndcg  precision    recall

Bias      200       10.0    0.367318  0.111554   0.151311  0.071239

Bias      200      100.0    0.377586  0.192387   0.059459  0.236567

AlgoClass features  nrecs  recip_rank      ndcg  precision    recall

BiasedMF  100        10.0    0.145303  0.043043   0.055902  0.022754

BiasedMF  100       100.0    0.163567  0.091865   0.029721  0.113809

ImplicitMF 20         10.0    0.447507  0.172096   0.215246  0.138612

ImplicitMF 20        100.0    0.447372  0.299152   0.095525  0.438936

Bias Test RMSE:  0.9418317549143906, 
Bias Test MAE:  0.7440482507386398,
Biasedmf Test RMSE:  0.8407684464066277,
Biasedmf Test MAE:  0.6449002519496809
### large dataset
AlgoClass damping   nrecs  recip_rank      ndcg  precision    recall

Bias      200      100.0     0.14971  0.082469   0.021677  0.142896

Bias      200       10.0    0.139054  0.040272   0.040102  0.031545

AlgoClass features  nrecs  recip_rank      ndcg  precision    recall

BiasedMF  200       10.0     0.028419  0.012167   0.009841  0.011476

BiasedMF  200       100.0    0.035671  0.032845   0.007739  0.060251

large Bias Test RMSE:  0.899735779778348,
large Bias Test MAE:  0.6941756862027079,
large Biasedmf Test RMSE:  0.8181884457208151,
large Biasedmf Test MAE:  0.6205422334228684.

Due to the time and memory constraint, I did not get the parameter-tuned result from implicitmf on full dataset. I have already tried something like cpus=5,nodes=4,hours=7,eval_jobs=28, but still can not get the implicitmf done. Meanwhile, the hpc is really crowded and I have to wait a long time for each ticket. 

## Discussion
What selection criteria (choice of validation metric) did you use for model selection (hyper-parameter tuning) and why?\
I use the damping for bias and features for biasedmf. implicitmf, because according to the experiments of lab1, they have a great impact on the evaluation results. Tuning these hyperparmeter will be effective and efficient. I use ndcg to measure the relevance since it is able to use the fact that some documents are “more” relevant than others. Highly relevant items should come before medium relevant items, which should come before non-relevant items.\
Would different model selection criteria have led to different outcomes?\
For most of them, No. I observe that the binary relevance metrics always have the same trend. For example, the larger the damping for Bias, the greater the relevance is. The Biasedmf get really close results for features=100 and 200, so it might cause different selection if we use different metric.\
Which method(s) perform best on the test for each of the metrics in question?\
MRR ImplicitMF\
ndcg ImplicitMF\
precision@10    ImplicitMF\
recall@10 ImplicitMF\
precision@100   ImplicitMF\
recall@100 ImplicitMF\
RMSE BiasedMF\
MAE BiasedMF
