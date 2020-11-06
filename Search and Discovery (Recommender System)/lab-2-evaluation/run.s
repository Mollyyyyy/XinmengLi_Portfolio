#!/bin/bash
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=1
#SBATCH --cpus-per-task=8
#SBATCH --time=02:00:00
#SBATCH --mem=20GB
#SBATCH --job-name=dsga3001-lab2
#SBATCH --mail-type=END
#SBATCH --mail-user=xl1575@nyu.edu
#SBATCH --output=slurm_lab2_%j.out

# Refer to https://sites.google.com/a/nyu.edu/nyu-hpc/documentation/prince/batch/submitting-jobs-with-sbatch
# for more information about the above options

# Remove all unused system modules
module purge

# Move into the directory that contains our code
SRCDIR=/scratch/xl1575/lab-2-evaluation-Mollyyyyy

# Activate the conda environment
source ~/.bashrc
conda activate dsga3001

# Execute the script
#stdbuf -o0 -e0 python rmsemae.py>rmse-lar
stdbuf -o0 -e0 python code_lar.py>failedres
#python count.py
# And we're done!
