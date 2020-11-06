
from lenskit.datasets import MovieLens
from lenskit.metrics.predict import rmse
from lenskit.batch import predict
from lenskit.algorithms.basic import Bias
from lenskit.algorithms.als import BiasedMF
from lenskit.util import rng
import numpy as np
from lenskit.crossfold import partition_users, SampleFrac
if __name__ == '__main__':
        N_SPLITS = 5
        FRAC_SPLIT = 0.2
        data = MovieLens('ml-latest-small')
        ratings = data.ratings[['user', 'item', 'rating']]
        n=1
        valli = []
        testli = []
        vallimf = []
        testlimf = []
        for train, test in partition_users(ratings,N_SPLITS,SampleFrac(FRAC_SPLIT),rng_spec=rng(0)):
                print("------------ trial ",n," -------------")
                n+=1
                rmli = []
                b =  list(partition_users(train,1,SampleFrac(FRAC_SPLIT),rng_spec=rng(0)))
                Train,val = b[0][0],b[0][1]
                print("####### Bias Model #######")
                dli = np.arange(10)
                biasli = []
                for d in dli:
                        bias = Bias(damping = d).fit(Train)
                        biasli.append(bias)
                        #preds = predict(bias, val)
                        preds = bias.predict(val[['user', 'item']])
                        rmli.append(rmse(preds, val['rating']))
                        #rmli.append(rmse(preds['prediction'], preds['rating']))
                print("rmse on validation for different damping values: ",rmli)
                i=rmli.index(min(rmli))
                print("when damping is ",dli[i],", we got the lowest rmse on validation set ",rmli[i])
                preds = predict(biasli[i], test)
                valli.append(rmli[i])
                t = rmse(preds['prediction'], preds['rating'])
                testli.append(t)
                print("rmse on test set is ",t)
                print("####### BiasMF Model #######")
                regli = [10**i for i in range(-3,2)]
                featli = np.arange(10,350,80)
                biasli = []
                rmli = []
                for f in featli:
                        for r in regli:
                                bias = BiasedMF(features=f,reg=r).fit(Train)
                                biasli.append(bias)
                                preds = bias.predict(val[['user', 'item']])
                                rmli.append(rmse(preds, val['rating']))
                                #preds = predict(bias, val)
                                #rmli.append(rmse(preds['prediction'], preds['rating']))
                print("rmse on validation for different feature and reg values: ",rmli)
                i=rmli.index(min(rmli))
                print("when feature is ",featli[i//len(regli)],", reg is ",regli[i%len(regli)],", we got the lowest rmse on validation set ",rmli[i])
                preds = predict(biasli[i], test)
                vallimf.append(rmli[i])
                t = rmse(preds['prediction'], preds['rating'])
                testlimf.append(t)
                print("rmse on test set is ",t)
        #print(valli,testli,vallimf,testlimf)
        print("Bias Model: Validation rmse across trials: Mean ",np.mean(np.array(valli)),"Standard Deviation ",np.std(np.array(valli)))
        print("Bias Model: Test rmse across trials: Mean ",np.mean(np.array(testli)),"Standard Deviation ",np.std(np.array(testli)))
        print("BiasedMF Model: Validation rmse across trials: Mean ",np.mean(np.array(vallimf)),"Standard Deviation ",np.std(np.array(vallimf)))
        print("BiasedMF Model: Test rmse across trials: Mean ",np.mean(np.array(testlimf)),"Standard Deviation ",np.std(np.array(testlimf)))
