%% ========================================== Harris & Feature matching
I1 = double(rgb2gray(imread('./img/building_1.jpg')));
[I11,I22] = size(I1);
response1 = harris_corners(I1); % get response map by Harris corner dector
keypoints1 = corner_peaks(response1, 45); % find keypoints(corners) according to the response
desc1 = describe_keypoints(I1, keypoints1); % describe the keypoints by your self-designed descriptor

I2 = double(rgb2gray(imread('./img/building_2.jpg')));
response2 = harris_corners(I2);
keypoints2 = corner_peaks(response2, 45);
desc2 = describe_keypoints(I2, keypoints2);

figure(1);
subplot(1,2,1), plot_keypoints(I1, keypoints1), title('Keypoints in image 1');
subplot(1,2,2), plot_keypoints(I2, keypoints2), title('Keypoints in image 2');
pause(0.5);

matches = match_descriptors(desc1, desc2); % match the keypoints accroding to your descriptors
points1 = keypoints1( matches(:, 1), :);
points2 = keypoints2( matches(:, 2), :);
match_plot(I1, I2, points1, points2); % visualize the match

%% ========================================== Try SURF features
I1 = rgb2gray(imread('./img/building_1.jpg'));
I2 = rgb2gray(imread('./img/building_2.jpg'));

% ------------------------------ Write your code here
% you need to find the corresponding points between two images by SURF feature, then visualize the matches
% you can use any built-in functions for this problem
points1 = detectSURFFeatures(I1);
points2 = detectSURFFeatures(I2);
[f1,vpts1] = extractFeatures(I1,points1);
[f2,vpts2] = extractFeatures(I2,points2);
indexPairs = matchFeatures(f1,f2) ;
matchedPoints1 = vpts1(indexPairs(:,1));
matchedPoints2 = vpts2(indexPairs(:,2));
figure; 
showMatchedFeatures(I1,I2,matchedPoints1,matchedPoints2);
legend('matched points 1','matched points 2');
% -------------- end of your code


%% ================================= Functions for Harris
function y = harris_corners(I) % get response map by Harris corner dector
    [n, m] = size(I);
    border=6;
    sigma=2;
    g = fspecial('gaussian', max(1,6*sigma), sigma); % Gaussian window function
    alpha = 0.04; % The constant in corner response function
    % ---------------------------------- Write your code here
    % you need to compute the response map according to corner response function
    % before that, you may need to get the image derivatives and the matrix M
    % since our local feature descriptor needs a patch, we only consider the area I(border+1:n-border,border+1:m-border)
    bI = I(border+1:n-border,border+1:m-border);
    %Sobel Derivative
    dx = [1 0 -1; 2 0 -2; 1 0 -1];
    dy = dx';
    Ix = conv2(bI,dx,'same');
    Iy = conv2(bI,dy,'same');
    Ix2 = conv2(Ix.^2, g, 'same');
    Iy2 = conv2(Iy.^2, g, 'same');
    Ixy = conv2(Ix.*Iy, g, 'same');
    y = (Ix2.*Iy2-Ixy.^2) - alpha*(Ix2+Iy2).^2;
    % -------------- end of your code
end

function y = corner_peaks(response, threshold) % find corners according to the response
    border=6;
    r=6;
    response=(1000/max(max(response)))*response;
    R=response;
    sze = 2*r+1; 
    MX = ordfilt2(R,sze^2,ones(sze));
    response = (R==MX)&(R>threshold); 
	R=R*0;
    R(5:size(response,1)-5,5:size(response,2)-5)=response(5:size(response,1)-5,5:size(response,2)-5);
	[r1,c1] = find(R);
    y = [r1+border+1,c1+border+1]; 
end

function y = plot_keypoints(I, PIP)
   Size_PI=size(PIP,1);
   for r=1: Size_PI
       I(PIP(r,1)-2:PIP(r,1)+2,PIP(r,2)-2)=255;
       I(PIP(r,1)-2:PIP(r,1)+2,PIP(r,2)+2)=255;
       I(PIP(r,1)-2,PIP(r,2)-2:PIP(r,2)+2)=255;
       I(PIP(r,1)+2,PIP(r,2)-2:PIP(r,2)+2)=255;
   end
   imshow(uint8(I));
   y = 0;
end

%% ================================= Functions for feature matching
function y = simple_descriptor(patch)
% design you own local feature descriptor
    % ---------------------------------- Write your code here
    % Reference : https://farshbaf.net/en/artificial-intelligence/blog/hog-matlab-implementation
    [m, m1] = size(patch);
    if m1 ~= m
        error('Patch is not a square matrix');
    end
    n= floor(m/3.5);
    border=2;
    patch=double(patch);
    k = 20;
    bI = patch(border+1:m-border,border+1:m-border);
    g = fspecial('gaussian',n);
    dx = [1 0 -1; 1 0 -1; 1 0 -1];
    dy = dx';
    Ix = conv2(bI,dx,'same');
    Iy = conv2(bI,dy,'same');
    magn = imfilter(sqrt(Ix.^2+Iy.^2),g);
    ang = atan2(Iy,Ix)*180/pi;
    [a,b] = size(ang);
    disp("a,b");
    disp(size(ang));
    a = floor(a/n);
    b = floor(b/n);
    bin = zeros(a,b,360/k);
    disp(size(bin));
    for i = 1:a
        for j = 1:b
            for bini=1:360/k
                A=zeros(n,n);
                for p = n*(i-1)+1:i*n
                    for q = n*(j-1)+1:j*n
                        if((ang(p,q)>=(bini-1)*k+1)&&(ang(p,q)<(bini)*k))
                            A(p-n*(i-1),q-n*(j-1))=1;
                        elseif(bini>1)
                            if((ang(p,q)>=(bini-2)*k+1+k/2)&&(ang(p,q)<(bini-1)*k))
                                A(p-n*(i-1),q-n*(j-1))=1-abs(ang(p,q)-(bini*k-k/2))/k;
                            end                
                        elseif(bini<360/k)
                            if((ang(p,q)>=(bini)*k+1)&&(ang(p,q)<(bini+1)*k-k/2))
                                A(p-n*(i-1),q-n*(j-1))=1-abs(ang(p,q)-(bini*k-k/2))/k;
                            end
                        end
                    end
                end
                bin(i,j,bini) = sum(sum( A.*magn(n*(i-1)+1:i*n,n*(j-1)+1:j*n)));
            end
        end
    end
    disp("bbb");
    disp(size(bin));
    block=zeros(a-1,b-1,4*360/k);
    for r=1:a-1
        for s=1:b-1
            vec=zeros(1,4*360/k);
            for i=1:2
                for j=1:2
                    num = ((i-1)*2+j);
                    vec((num-1)*(360/k)+1:(num)*(360/k)) = permute(bin(r-1+i,s-1+j,:),[3,2,1]);
                end
            end
            %norm=vec./sum(vec);
            norm=vec ./ sum(vec);
            norm(norm(:)>0.2)=0.2;
            norm=norm ./ sum(norm);
            block(r,s,:)=norm;
        end
    end
    % Convert block feature to vector
    blocks=zeros(1,size(block,1)*size(block,2)*size(block,3));
    for i=1:size(block,1)
        for j=1:size(block,2)
            for o=1:size(block,3)
            blocks((i-1)*size(block,2)*size(block,3)+(j-1)*size(block,3)+o)=block(i,j,o);
            end
        end
    end
    y = blocks;
    % -------------- end of your code
end

function ret = describe_keypoints(I, keypoints)
% using the above simple_descriptor() to describle the keypoints
% you may need to find a patch for each keypoint, then call the simple_descriptor()
% your output size should be [num_of_keypoints, feature_size]
    patch_size = 12;
    % ---------------------------------- Write your code here
    [a,b] = size(I);
    patch_size = floor(a/8);
    [m,n] = size(keypoints);
    loc = 0;
    row = keypoints(1,1);
    col = keypoints(1,2);
    patch = I(row:row+patch_size-1,col:col+patch_size-1);
    disp("patch size");
    disp(size(patch));
    re = simple_descriptor(patch);
    %[re,hogVisualization] = extractHOGFeatures(patch);
    disp("descriptor");
    disp(size(re));
    for i=2:m
        row = keypoints(i,1);
        col = keypoints(i,2);
        % Let the keypoint be at the center of patch
        %if row-5 > 0 && col-5 > 0 && row+6 <= a &&col+6 <= b
         %   loc = -6;
        %end
        for pati = 1:patch_size
            for patj = 1:patch_size
                if pati+loc+row <= a && patj+loc+col <=b
                    patch(pati,patj)=I(pati+loc+row,patj+loc+col);
                end
            end    
        end
        re = [re;simple_descriptor(patch)];
        %re = [re;extractHOGFeatures(patch)];
    end
    ret = re;
    disp('ret');
    disp(size(ret));
    % -------------- end of your code
end

function y = match_descriptors(desc1, desc2)
% match the keypoints according to your descriptors
% you may need to adjust the threshold according to your distance function
% your output size should be [num_of_matches, 2]
    % ---------------------------------- Write your code here
    [a,b] = size(desc1);
    [c,d] = size(desc2);
    disp('a');
    disp(a);
    disp(b);
    disp('c');
    disp(c);
    disp(d);
    disp('dd');
    mat = [-1,-1];
    for i=1:a
        d = zeros(c,1);
        d1 = desc1(i,:);
        for j=1:c
            d2 = desc2(j,:);
            %d(j) = acos((d1'*d2));
            d(j) = sqrt(sum((d1 - d2).^2));
        end
        [sor,ind] = sort(d);
        if sor(1)/sor(2) <= 0.8 || sor(2)==0
            if isequal(mat,[-1,-1]) == 1
                mat=[i,ind(1)];
            else
                mat = [mat;i,ind(1)];
            end
        end
    end
    y = mat;
    % -------------- end of your code
end

function h = match_plot(img1,img2,points1,points2)
    h = figure;
    colormap = {'b','r','m','y','g','c'};
    height = max(size(img1,1),size(img2,1));
    match_img = zeros(height, size(img1,2)+size(img2,2), size(img2,3));
    match_img(1:size(img1,1),1:size(img1,2),:) = img1;
    match_img(1:size(img2,1),size(img1,2)+1:end,:) = img2;
    imshow(uint8(match_img));
    hold on;
    for i=1:size(points1,1)
        plot([points1(i,2) points2(i,2)+size(img1,2)],[points1(i,1) points2(i,1)],colormap{mod(i,6)+1});
    end
    title('Correspondence');
    hold off;
end
