%% ======================== Images
% Open image and convert it to grayscale
img_1 = imread('./img/dog.jpg');
img_1 = rgb2gray(img_1);
% Another very simple image
img_2 = ones(10);
img_2(2:6,2:6)=-1;
% Show image
img = img_1; % you can choose img_1 or img_2 here, or other images you like
figure(1);
imshow(img);
title('Original image');
pause(0.1);

%% ======================== Filters
filter_1= [1 0 -1;
 2 0 -2;
 1 0 -1];
filter_2 = [
1 1 -1 1 2 4;
1 1 1 1 3 2;
1 1 -1 1 2 5;
1 1 1 1 2 1;
] * (1/16);

filter = filter_1; % you can choose filter_1 or filter_2 here, or other filters you like

%% ======================== Convolution
my_output1 = my_conv2d(img, filter);
expected_outoput1 = conv2(img, filter, 'same');
figure(2);
subplot(1,2,1), imshow(uint8(my_output1)), title('My output convolution');
subplot(1,2,2), imshow(uint8(expected_outoput1)), title('Expected output convolution');
pause(0.1);

% ----- Your my_conv2d function should get the same results with the built-in function
if max(my_output1 - expected_outoput1) < 1e-6
    fprintf('Your conv function is correct.\n');
else
    fprintf('Your conv function is not correct.\n');
end

%% ======================== Correlation
my_output = my_corr2d(img, filter);
expected_outoput = filter2(filter, img);
figure(3);
subplot(1,2,1), imshow(uint8(my_output)), title('My output correlation');
subplot(1,2,2), imshow(uint8(expected_outoput)), title('Expected output correlation');
pause(0.1);

% ----- Your my_corr2d function should get the same results with the built-in function
if max(my_output - expected_outoput) < 1e-6
    fprintf('Your corr function is correct.\n');
else
    fprintf('Your corr function is not correct.\n');
end

%% =========================
function y = my_conv2d(I,f)
% I is the image, f is the filter
    I = double(I);
    % ------------------------ Write your code here
    f = double(f);
    [a,b]=size(I);
    [c,d]=size(f);
    mic = floor((c/2))+1;
    mid = floor((d/2))+1;
    y = zeros(a,b);
    for m=1:a
        for n=1:b
            for k=1:c
                for l=1:d
                    if m+mic-k>0 && n+mid-l>0 && n+mid-l <=b && m+mic-k <=a 
                        y(m,n)=y(m,n)+f(k,l)*I(m+mic-k,n+mid-l);
                    end
                end
            end
        end
    end
    % -------------- end of your code
end

function y = my_corr2d(I,f)
% I is the image, f is the filter
    I = double(I);
    % ------------------------ Write your code here
    f = rot90(f,2);
    f = double(f);
    [a,b]=size(I);
    [c,d]=size(f);
    mic = floor((c/2))+1;
    mid = floor((d/2))+1;
    y = zeros(a,b);
    for m=1:a
        for n=1:b
            for k=1:c
                for l=1:d
                    %if m>k+mic && n>l+mid && m-k+mic <=a && n-l+mid <=b
                    %    y(m,n)=y(m,n)+f(k,l)*I(m-k+mic,n-l+mid);
                    %end
                    if m+mic-k>0 && n+mid-l>0 && n+mid-l <=b && m+mic-k <=a 
                        y(m,n)=y(m,n)+f(k,l)*I(m+mic-k,n+mid-l);
                    end
                end
            end
        end
    end
    % -------------- end of your code
end