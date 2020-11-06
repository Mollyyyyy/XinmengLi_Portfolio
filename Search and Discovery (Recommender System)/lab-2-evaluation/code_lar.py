from lenskit.batch import MultiEval
from lenskit.datasets import MovieLens
from lenskit.metrics.predict import rmse,mae
from lenskit.batch import predict
from lenskit import topn
from lenskit.algorithms.basic import Bias
from lenskit.algorithms.als import BiasedMF
from lenskit.algorithms.als import ImplicitMF
import time
import numpy as np
from lenskit.crossfold import partition_users, SampleFrac,TTPair
import pandas as pd


def compare(folder,pair_binary,attr):
	runs = pd.read_csv(folder+'/runs.csv')
	runs.set_index('RunId', inplace=True)
	recs = pd.read_parquet(folder+'/recommendations.parquet')
	if 'exp' in folder:
		bi = (recs['score']>=4).astype(int)
		recs['score'] = bi
	truth = pd.concat((p.test for p in pair_binary), ignore_index=True)
	rla = topn.RecListAnalysis()
	rla.add_metric(topn.recip_rank)
	rla.add_metric(topn.ndcg)
	rla.add_metric(topn.precision)
	rla.add_metric(topn.recall)
	raw_metric = rla.compute(recs, truth)
	metric = raw_metric.join(runs[['AlgoClass']+attr], on='RunId')
	res = metric.fillna(0).groupby(['AlgoClass']+attr).mean()
	return res

if __name__ == '__main__':

	ratings  = MovieLens('/scratch/courses/DS-GA.3001-016-2020/ml-latest').ratings[['user', 'item', 'rating']]
	ratings['binary'] = (ratings['rating']>=4).astype(int)
	FRAC_SPLIT = 0.2
	pair =  list(partition_users(ratings,1,SampleFrac(FRAC_SPLIT,rng_spec=0),rng_spec=0))
	train_val,test = pair[0][0][['user', 'item','rating']],pair[0][1][['user', 'item','rating']]
	tratestpair = [TTPair(train_val,test)]
	train_val_binary,test_binary = pair[0][0][['user', 'item']],pair[0][1][['user', 'item']]
	train_val_binary['rating'] = pair[0][0]['binary']
	test_binary['rating'] = pair[0][1]['binary']
	tratestpair_binary = [TTPair(train_val_binary,test_binary)]
	travalpair = list(partition_users(train_val,1,SampleFrac(FRAC_SPLIT,rng_spec=0),rng_spec=0))
	travalpair_binary = list(partition_users(train_val_binary,1,SampleFrac(FRAC_SPLIT,rng_spec=0),rng_spec=0))
	names = ['Bias','BiasedMF','ImplicitMF']

	for r in [100]:
		
		#Bias
		eval = MultiEval('expval-bias-lar100', recommend=r)
		eval.add_datasets(travalpair, name='ML')
		eval.add_algorithms([Bias(damping=d) for d in [10,100,200]],attrs=['damping'], name=names[0])
		eval.run()

		res = compare('expval-bias-lar100',travalpair_binary,['damping'])
		print("--- Explicit Validation Performance ---")
		print(res)
		
		op1 = res['ndcg'].loc[names[0]].idxmax()
			
			#BiasedMF
		eval = MultiEval('expval-biasedmf-lar100', recommend=r)
		eval.add_datasets(travalpair, name='ML')
		eval.add_algorithms([BiasedMF(features=f) for f in [20,100,200]], attrs=['features'], name=names[1])
		eval.run()
	
		res = compare('expval-biasedmf-lar100',travalpair_binary,['features'])
		print(res)
		
		op2 = res['ndcg'].loc[names[1]].idxmax()

			#ImplicitMF
		eval = MultiEval('impval-lar100', recommend=r)
		eval.add_datasets(travalpair_binary, name='ML-binary')
		eval.add_algorithms([ImplicitMF(features=f) for f in [20,100,200]], attrs=['features'], name=names[2])
		eval.run()
	
		res = compare('impval-lar100',travalpair_binary,['features'])
		print("--- Implicit Validation Performance ---")
		print(res)  
	
		print(names[0], "Optimal damping: ", op1)
		print(names[1], "Optimal features: ", op2)
		op3 = res['ndcg'].loc[names[2]].idxmax()
		print(names[2], "Optimal features: ",op3)
		#Test
		#Bias

		eval = MultiEval('exptest-bias-lar100', recommend=r)
		eval.add_datasets(tratestpair, name='ML')
		eval.add_algorithms([Bias(damping=int(op1))],attrs=['damping'], name=names[0])
		eval.run()

		res = compare('exptest-bias-lar100',tratestpair_binary,['damping'])
		print("--- Test Performance ---")
		print(res)
		
		#BiasedMF
		eval = MultiEval('exptest-biasedmf-lar100', recommend=r)
		eval.add_datasets(tratestpair, name='ML')
		eval.add_algorithms([BiasedMF(features=int(op2))], attrs=['features'], name=names[1])  
		eval.run()
		res = compare('exptest-biasedmf-lar100',tratestpair_binary,['features'])
		print(res)
		op3 = 20
		#ImplicitMF
		eval = MultiEval('imptest-lar100', recommend=r)
		eval.add_datasets(tratestpair_binary, name='ML-binary')
		eval.add_algorithms([ImplicitMF(features=int(op3))], attrs=['features'], name=names[2])
		eval.run()
		
		res = compare('imptest-lar100',tratestpair_binary,['features'])
		print(res)
		
