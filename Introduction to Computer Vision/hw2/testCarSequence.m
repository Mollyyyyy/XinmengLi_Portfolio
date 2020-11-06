%*************************************************************************%
% Function:    testCarSequence                                            % 
% Description: Run the iterative KLT Tracker to track the speeding car    %
%              in the given video                                         %
%                                                                         %
%              box - a n*4 matrix, n is number of frames and each row     %
%                    in the matrix contains 4 numbers: [x1,y1,x2,y2]      %
%                    which indicates the coordinates of top-left and      %
%                    bottom-right corners of the tracked rectangles       %
%                                                                         %
%*************************************************************************%

% Loads the video sequence
load('carSequence.mat') ;
N_frames = size(sequence, 4) ;

% Initialize box
box = zeros(N_frames, 4) ;

% Initial position of the rectangle
rect = [328, 213, 419, 265]' ;
%rect = [328,213,333,215]';
box(1,:) = rect';

% ------------------------------ Write your code here
% Compute the width and height of the rectangle
height = rect(4)-rect(2)+1;
width = rect(3)-rect(1)+1;

% Get the first frame
It = sequence(:, :, :, 1);
It_color = It ;
It = rgb2gray (It) ;
It = im2double (It)  ;

% Draw the initial rectangle
figure ;
imshow(It_color) ;

rectangle('Position',[rect(1), rect(2), width, height], 'LineWidth',1.5,'edgecolor','y') ;
drawnow() ;    
nor=[0 0];
% Track the car from the second frame
for i = 2:N_frames
    
    % Read next frame in the sequence 
    It1 = sequence(:, :, :, i);
    It_color = It1;
    It1 = rgb2gray (It1);
    It1 = im2double (It1);
    % call iterative_KLT to compute [u, v]
    [u, v,delta_p,Ixx,Ix1,Inte,Inte1,X,X1,Y,Y1,A,B,ATB,sumA,nor] = iterative_KLT(It,It1,rect,nor);
    % Update rectangle's position
    rect = [round(rect(1)+v),round(rect(2)+u),round(rect(1)+v)+width-1,round(rect(2)+u)+height-1]';
    % Record rectangle in box
    box(i,:) = rect';
    % Draw the rectangle
    imshow( It_color ) ;
    rectangle('Position',[rect(1), rect(2), width, height], 'LineWidth',1.5,'edgecolor','y') ;
    drawnow() ;
    
    % Go to the next frame
    It = It1;
end
% -------------- end of your code

% Save box in carPosition.mat
save('carPosition.mat', 'box') ;


