%% ============================ Prepare data and model
[imgDataTrain, labelsTrain, imgDataTest, labelsTest] = prepareData;
load MNISTModel
net.Layers
% if you successfully installed the Deep Learning toolbox, you should see the architecture of the CNN in the command window

%% ============================ Try prediction
% Just try the following codes, they predict the class of an image
randIndx = randi(numel(labelsTest));
img = imgDataTest(:,:,1,randIndx);
actualLabel = labelsTest(randIndx);
predictedLabel = net.classify(img);
figure(1);
imshow(img);
title(['Predicted: ' char(predictedLabel) ', Actual: ' char(actualLabel)])

%% ============================= Visualize filters in conv1
% ------------------------------ Write your code here
% you need to visualize the filters in conv1, then explain the effect of two filters in your report
% fill the below blank
% you can write other codes before that
layer = 2;
disp(net.Layers(layer).Name);
disp(net.Layers(layer).Weights(:,:,1,1));
figure(2);
for i = 1:16
    subplot(4,4,i),imshow(mat2gray(conv2(img,net.Layers(layer).Weights(:,:,1,i),'same')));
    title(['filter ' num2str(i)]);
end
% -------------- end of your code

%% ============================= Network activations of conv1
% ------------------------------ Write your code here
% select 5 samples in imgDataTest, and visualize the activations of conv1 based on the above 2 filters
%filter1=net.Layers(2).Weights(:,:,1,2);
%filter2=net.Layers(2).Weights(:,:,1,6);
figure(3);
for i = 1:5
    randIndx = randi(numel(labelsTest));
    im = imgDataTest(:,:,1,randIndx);
    act1 = activations(net,im,'conv_1');
    sz = size(act1);
    act1 = reshape(act1,[sz(1) sz(2) 1 sz(3)]);
    subplot(5,2,(i-1)*2+1), imshow(mat2gray(act1(:,:,1,2)));
    subplot(5,2,(i-1)*2+2), imshow(mat2gray(act1(:,:,1,6)));
end
% -------------- end of your code

%% ============================= Image retrieval based on conv3
k = 5;
search_num = 1000;
% ------------------ Write your code here
% select 3 images in imgDataTest, for each image
% find the 5 most similar images in the first 1000 images of imgDataTrain
% use Euclidean distance on conv3's activations
imtra = imgDataTrain(:,:,1,1:search_num);
disp(size(imtra));
act = activations(net,imtra,'conv_3');
sz = size(act);
act = reshape(act,[sz(1) sz(2) 1 sz(3) sz(4)]);
disp(size(act))
for i = 1:3
    randIndx = randi(numel(labelsTest));
    im = imgDataTest(:,:,1,randIndx);
    act2 = activations(net,im,'conv_3');
    sz2 = size(act2);
    act2 = reshape(act2,[sz2(1) sz2(2) 1 sz2(3)]);
    disp(size(act2));
    act2v = reshape(act2,[],1);
    dis = zeros(search_num,1);
    for j = 1:search_num
        actv = reshape(act(:,:,1,:,j),[],1);
        dis(j) = sqrt(sum((actv - act2v).^2));
    end
    [sor,index] = sort(dis);
    figure(3+i);
    subplot(1,k+1,1), imshow(im);
    title('Original Image');
    for q = 1:k
        subplot(1,k+1,q+1), imshow(imgDataTrain(:,:,1,index(q)));
    end
end
% -------------- end of your code
