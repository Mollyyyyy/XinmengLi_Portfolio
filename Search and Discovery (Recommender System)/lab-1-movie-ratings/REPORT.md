# Lab 1 - Report

Name: Xinmeng Li\
NetID: xl1575

## Description of experiment
### What hyper-parameters did you try for each model?  
Bias Model: damping \
BiasedMF Model: feature and reg
### How did you pick the range of values to test?
The damping parameter range is the integer 0-9 generated by np.arange(10). \
The reg range is 0.001, 0.01, 0.1, 1, 10\
The feature range at first was small, i.e. feature<=50. I observed that in this case, the model performs best when feature=50 in all of 5 trials (Please refer to [link](https://github.com/NYU-Search-Discovery/lab-1-movie-ratings-Mollyyyyy/blob/master/out) for the result). Therefore, I increase the upperbound of feature to be 330. The final range for feature is 10,  90, 170, 250, 330 (Please refer to [link](https://github.com/NYU-Search-Discovery/lab-1-movie-ratings-Mollyyyyy/blob/master/out2) for the result).
### How many users, items, and ratings were used (on average) by your train, validation, and test splits?
On average, I got\
Train: 610 users, 8834 items, 77447 ratings\
Validation: 610 users, 5043 items, 19356 ratings\
Test: 122 users, 2188 items, 4033 ratings
## Results
I compared the validation rmse on all of the candidate values of hyperparameters, then calculate the test rmse on the best setting, i.e. hyperparameters with the lowest validation rmse. Finally, I calculate the mean and standard deviation of the test rmse across 5 trials for Bias and BiasedMF respectively.
### Here is the result for above procedure.
------------ trial  1  ------------- \
####### Bias Model ####### \
rmse on validation for different damping values:  [0.9009, 0.8806, 0.8761, 0.8747, 0.8744, 0.8746, 0.8751, 0.8757, 0.8764, 0.8772] \
when damping is  4 , we got the lowest rmse on validation set  0.874376412913912 \
rmse on test set is  0.8838562261582306 \
####### BiasMF Model ####### \
rmse on validation for different feature and reg values:  [1.3804, 1.0929, 0.866, 0.8657, 0.8657, 1.026, 0.8737, 0.8447, 0.8657, 0.8657, 0.9645, 0.863, 0.8445, 0.8657, 0.8657, 0.9688, 0.8625, 0.8445, 0.8657, 0.8657, 0.9637, 0.8626, 0.8445, 0.8657, 0.8657]\
when feature is  330 , reg is  0.1 , we got the lowest rmse on validation set  0.8445021707277439\
rmse on test set is  0.8686121375283328 \
------------ trial  2  -------------\
####### Bias Model #######\
rmse on validation for different damping values:  [0.8908, 0.8685, 0.8637, 0.8623, 0.8622, 0.8626, 0.8633, 0.8642, 0.8651, 0.8661]\
when damping is  4 , we got the lowest rmse on validation set  0.8621803094740349\
rmse on test set is  0.8692495273428419\
####### BiasMF Model #######\
rmse on validation for different feature and reg values:  [1.4105, 1.086, 0.8563, 0.8569, 0.8569, 1.0418, 0.8684, 0.837, 0.8569, 0.8569, 0.9701, 0.8582, 0.837, 0.8569, 0.8569, 0.9545, 0.8571, 0.8368, 0.8569, 0.8569, 0.9657, 0.8574, 0.8369, 0.8569, 0.8569]\
when feature is  250 , reg is  0.1 , we got the lowest rmse on validation set  0.8368480533614077\
rmse on test set is  0.836238567176803\
------------ trial  3  -------------\
####### Bias Model #######\
rmse on validation for different damping values:  [0.8905, 0.8729, 0.8694, 0.8687, 0.8689, 0.8695, 0.8704, 0.8713, 0.8723, 0.8733]\
when damping is  3 , we got the lowest rmse on validation set  0.8686585557832485\
rmse on test set is  0.823210334795582\
####### BiasMF Model #######\
rmse on validation for different feature and reg values:  [1.3899, 1.0902, 0.8623, 0.8627, 0.8627, 1.0438, 0.873, 0.8451, 0.8627, 0.8627, 0.9717, 0.8639, 0.8449, 0.8627, 0.8627, 0.9664, 0.863, 0.8448, 0.8627, 0.8627, 0.9643, 0.8639, 0.8449, 0.8627, 0.8627]\
when feature is  250 , reg is  0.1 , we got the lowest rmse on validation set  0.8448350315057915\
rmse on test set is  0.8014054942598158\
------------ trial  4  -------------\
####### Bias Model #######\
rmse on validation for different damping values:  [0.9001, 0.8772, 0.8723, 0.871, 0.871, 0.8715, 0.8722, 0.8731, 0.8741, 0.8751]\
when damping is  4 , we got the lowest rmse on validation set  0.8709845981587834\
rmse on test set is  0.8839667935231663\
####### BiasMF Model #######\
rmse on validation for different feature and reg values: [1.369, 1.0889, 0.8665, 0.865, 0.865, 1.0343, 0.8757, 0.8473, 0.865, 0.865, 0.9677, 0.8652, 0.8472, 0.865, 0.865, 0.9664, 0.8634, 0.8471, 0.865, 0.865, 0.9667, 0.864, 0.8471, 0.865, 0.865]\
when feature is  330 , reg is  0.1 , we got the lowest rmse on validation set  0.8470533216594587\
rmse on test set is  0.8639822535418681\
------------ trial  5  -------------\
####### Bias Model #######\
rmse on validation for different damping values:  [0.8917, 0.8737, 0.8701, 0.8694, 0.8696, 0.8702, 0.8711, 0.872, 0.873, 0.874]\
when damping is  3 , we got the lowest rmse on validation set  0.869363894630914\
rmse on test set is  0.8678330766090511\
####### BiasMF Model #######\
rmse on validation for different feature and reg values: [1.4163, 1.1034, 0.8647, 0.8634, 0.8634, 1.0331, 0.8759, 0.8461, 0.8634, 0.8634, 0.9673, 0.8663, 0.8462, 0.8634, 0.8634, 0.9547, 0.8658, 0.8461, 0.8634, 0.8634, 0.962, 0.8652, 0.8461, 0.8634, 0.8634]\
when feature is  330 , reg is  0.1 , we got the lowest rmse on validation set  0.8460751028828488\
rmse on test set is  0.8468652750730186

Bias Model: Validation rmse across trials: Mean  0.8691127541921786 Standard Deviation  0.003988434792220027\
Bias Model: Test rmse across trials: Mean  0.8656231916857744 Standard Deviation  0.02229715500104139\
BiasedMF Model: Validation rmse across trials: Mean  0.8438627360274502 Standard Deviation  0.0036231039708091647\
BiasedMF Model: Test rmse across trials: Mean  0.8434207455159676 Standard Deviation  0.024025451718968183
## Discussion
We observe that \
The averagae rmse on test and validation set are close. \
The rmse on test set has a larger standard deviation than the validation set.
### What was the best (average test) RMSE achieved in your experiments, and for what setting of hyper-parameters?
The best mean of test rmse is 0.8434207455159676 achieved by BiasedMF Model, with reg=0.1, feature=330 or 250. \
For BiasedMF model, the lowest test rmse among 5 trials is 0.8014054942598158 with reg=0.1, feature= 250.\
For Bias Model, the lowest test rmse among 5 trials is 0.823210334795582 with damping=3.