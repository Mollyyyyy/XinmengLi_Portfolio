from lenskit.metrics.predict import rmse,mae
import pandas as pd
if __name__ == '__main__':
        pre = pd.read_parquet('expval-bias-lar100/predictions.parquet')
        print("large Bias Validation RMSE: ",rmse(pre['prediction'], pre['rating']))
        print("large Bias Validation MAE: ",mae(pre['prediction'], pre['rating']))
        pre = pd.read_parquet('exptest-bias-lar100/predictions.parquet')
        print("large Bias Test RMSE: ",rmse(pre['prediction'], pre['rating']))
        print("large Bias Test MAE: ",mae(pre['prediction'], pre['rating']))
        pre = pd.read_parquet('expval-biasedmf-lar100/predictions.parquet')
        print("100 large Biasedmf Validation RMSE: ",rmse(pre['prediction'], pre['rating']))
        print("100 large Biasedmf Validation MAE: ",mae(pre['prediction'], pre['rating']))
        pre = pd.read_parquet('expval-biasedmf-lar/predictions.parquet')
        print("10 large Biasedmf Validation RMSE: ",rmse(pre['prediction'], pre['rating']))
        print("10 large Biasedmf Validation MAE: ",mae(pre['prediction'], pre['rating']))
