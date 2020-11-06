%*************************************************************************%
% Script for Q2.2 of homework 3                                           %
%*************************************************************************%

% Read the two images and convert them to double
I1 = imread('temple/im1.png') ;
I1 = im2double(I1) ;
I2 = imread('temple/im2.png') ;
I2 = im2double(I2) ;

% Get the scalar M, the largest image dimension
M = max(size(I1, 1), size(I1, 2)) ;

% Load the correspondences for the eight point algorithm
load('temple/some_corresp.mat') ;

% Compute the fundamental matrix using the eight point algorithm
F = eightpoint(pts1, pts2, M) ;

% ------- YOUR CODE HERE

% 1. Load K1 and K2
load('C:\Users\xinme\Desktop\152\A92092169_hw3\hw3_release\temple\intrinsics.mat');
% 2. Find M2 and M1
M2 = camera2(F,K1,K2,pts1,pts2);
M2 = K2*M2;
I1 = eye(3,4);
M1 = K1*I1;
% 3. Load the correspondences for 3D visualization
load('many_corresp.mat')
% 4. Get 3D points given 2D point correspondences
p1 = zeros(size(x1,1),2);
p2 = zeros(size(x2,1),2);
p1(:,1) = x1;
p1(:,2) = y1;
p2(:,1) = x2;
p2(:,2) = y2;
P = triangulate( M1, p1, M2, p2 );
% 5. Plot 3D points
scatter3(P(:,1),P(:,2),P(:,3));
% ------- END OF YOUR CODE
