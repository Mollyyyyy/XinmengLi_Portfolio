#!/bin/bash
#SBATCH --nodes=12
#SBATCH --ntasks-per-node=3
#SBATCH --cpus-per-task=1
#SBATCH --time=1:00:00
#SBATCH --mem=2GB
#SBATCH --job-name=dsga3001-lab1
#SBATCH --mail-type=END
#SBATCH --mail-user=xl1575@nyu.edu
#SBATCH --output=slurm_lab0_%j.out

# Refer to https://sites.google.com/a/nyu.edu/nyu-hpc/documentation/prince/batch/submitting-jobs-with-sbatch
# for more information about the above options

# Remove all unused system modules
module purge

# Move into the directory that contains our code
SRCDIR=/scratch/xl1575/lab-1-movie-ratings-Mollyyyyy

# Activate the conda environment
source ~/.bashrc
conda activate dsga3001

# Execute the script
stdbuf -o0 -e0 python code.py>out3
# And we're done!
