%*************************************************************************%
% Function:    epipolarCorrespondence                                     % 
% Description: Slide a window along the epipolar line in im2 and find     %
%              the one that matches most closely to the window around     %
%              the point in im1                                           %  
%                                                                         %
%              Input:  x1 and y1 - The x and y coordinates of a pixel     %
%                                  on im1                                 %
%                                                                         %
%                      F - The 3*3 fundamental matrix                     %
%                                                                         %
%              Output: x2 and y2 - The x and y coordinates of the         %
%                                  corresponding pixel on im1             %
%*************************************************************************%

function [x2, y2] = epipolarCorrespondence(im1, im2, F, x1, y1)

% ------- YOUR CODE HERE
% (note: the code and comments below are provided for reference only)
% (feel free to change or add code if necessary)

% Convert RGB to gray
im1 = rgb2gray(im1);
im2 = rgb2gray(im2);

% Compute the size of the window depending on the value of the gaussian sigma 
sigma = 6;
window_size = 3*ceil(2*sigma)+1;
x1 = round(x1);
y1 = round(y1);
% Get the patch in image1
patch1 = getWindow(x1, y1, window_size, im1);

% Compute the epipolar line
l = F*[x1;y1;1];
disp(l);
% Get the range of x, y in im2
[ysize,xsize] = size(im2);
xrange = x1-60:0.05:x1+60;
% Get the points along the line
x2 = [];
y2 = [];
for i=xrange
    j = (-l(3)-l(1)*i)/l(2);
    if j <= y1+20 && j >= y1-20
        y2 = [y2;j];
        x2 = [x2;i];
    end
end
disp(y2);
% Limit the searching range to a neighborhood

% Get the points that are close to the one in image 1

% Find the correspondence in im2
n = size(x2,1);
disp(n);
mind = 100000;
xmin = 0.0;
ymin = 0.0;
for i=1:n
    patch2 = getWindow( x2(i), y2(i), window_size, im2);
    d = computeDifference(patch1, patch2, window_size, sigma) ;
    if d < mind
        mind = d;
        xmin = x2(i);
        ymin = y2(i);
    end
end
x2 = xmin;
y2 = ymin;
% ------- END OF YOUR CODE

end

%*************************************************************************%
% Function:    getWindow                                                  % 
% Description: Get the intensity values of the window in the given image  %  
%                                                                         %
%              Input:  x and y - The coordinates of the window center     %
%                      n_row, n_col - The number of rows and columns of   %
%                                     the image                           %
%                      sizeW - The width of the window                    %
%                                                                         %
%              Output: patch                                              %
%*************************************************************************%
function [patch] = getWindow(x, y, sizeW, image)

% ------- YOUR CODE HERE
% (note: the code and comments below are provided for reference only)
% (feel free to change or add code if necessary)

% Compute the boundaries
% Get patch
w = round((sizeW-1)/2);
xx = floor(x)-w:floor(x)+w+1;
yy = floor(y)-w:floor(y)+w+1;
xxx = x-w:x+w+1;
yyy = y-w:y+w+1;
[X,Y] = meshgrid(xx,yy);
[Xq,Yq] = meshgrid(xxx,yyy);
patch = interp2(X,Y,image(yy,xx),Xq,Yq);
patch = patch(1:end-1,1:end-1);
% ------- END OF YOUR CODE

end

%*************************************************************************%
% Function:    computeDifference                                          % 
% Description: compute difference between the two windows                 %  
%                                                                         %
%              Input:  patch1 - patch from image1                         %
%                      patch2 - patch from image2                         %
%                      sizeW - The width of the window                    %
%                      sigma                                              %
%                                                                         %
%              Output: d - distance                                       %
%*************************************************************************%
function [d] = computeDifference(patch1, patch2, sizeW, sigma)
%fspecial
% ------- YOUR CODE HERE
% (note: the code and comments below are provided for reference only)
% (feel free to change or add code if necessary)
%h = fspecial('gaussian', [sizeW sizeW], sigma);
diff = patch1-patch2;
%d = norm(h.*diff);
d = sqrt(sum(diff.^2,'all'));
% ------- END OF YOUR CODE

end
