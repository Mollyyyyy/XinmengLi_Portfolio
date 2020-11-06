%***************************************************************************%
% Function:    iterative_KLT                                                % 
% Description: iterative KLT Tracker                                        %
%              Computes the optimal local motion (u,v) from frame It to     %
%              frame It+1 that minimizes the pixel squared difference       %
%              in the two rectangles                                        %
%              Assumes that the rectangle undergoes constant motion         %
%              in a small region rect                                       %
%                                                                           %
%              Input:  It - the image frame It with pixel values in [0,1]   % 
%                      It1 - the image frame It+1 with pixel values in [0,1]%
%                      rect - 4*1 vector [x1, y1, x2, y2]'                  %
%                             where (x1, y1) is the top-left corner         %
%                             (x2, y2) is the bottom-right corner           %
%                             x1,y1,x2,y2 may not be integers               %
%                             The rectangle is inclusive, i.e. it includes  %
%                             all the four corners.                         %
%                                                                           %
%              Output: (u,v) - the movement of rect from It to It1          % 
%                                                                           %
%***************************************************************************%

function [u, v,delta_p,Ixx,Ix1,Inte,Inte1,X,X1,Y,Y1,A,B,ATB,sumA,nor] = iterative_KLT(It, It1, rect,nor)

% ------------------------------ Write your code here
y1 = rect(1);
y2 = rect(3);
x1 = rect(2);
x2 = rect(4);
% Initialize u and v
u = 4;
v = -2;
%Pad the rectangle for interp2
uu = 1;
vv = 1;
% Get the size of image It1
[a,b] = size(It);
% Compute gradient of It+1
[Ix,Iy] = imgradientxy(It1);
x11=x1-uu;
x22=x2+uu ;
y11 = y1-vv ;
y22 = y2+vv ;
Ixx = Ix(x11:x22,y11:y22);
Iyy = Iy(x11:x22,y11:y22);
[X,Y] = meshgrid(y11:y22,(x11:x22)');
% Compute the intensity of the rectangle on It
Inte = It(x11:x22,y11:y22);
%disp('IIIINNNNNNTTTEEE');
%disp(size(Inte));
[R,Q] = size(It(x1:x2,y1:y2));
% Initialize the norm of delta_p 
delta_p_norm = 10;
X1=0;
Y1=0;
i = 1;
Inte1 = 0;
Ix1 = 0;
Iy1 = 0;
A = 0;
B =0;
delta_p = 0;
ATB=0;
nor1 = [];
nor2 = [];

% Iteratively compute u, v
while(delta_p_norm > 0.001)
    % Compute the new rectangle based on current u, v 
    % If you encountered numerical issues here, you can try to add a very small value to each coordinate, e.g. 1e-6 
    [X1,Y1] = meshgrid(y11+v:y22+v,(x11+u:x22+u)');
    Inte1 = interp2(X,Y,Inte,X1,Y1);
    Inte1(isnan(Inte1))=0;
    Ix1 = interp2(X,Y,Ixx,X1,Y1);
    Ix1(isnan(Ix1))=0;
    Iy1 = interp2(X,Y,Iyy,X1,Y1);
    Iy1(isnan(Iy1))=0;
    % Check that whether the new rectangle is out of the image boundary
    if ( x11+u<1||y11+v<1||x22+u>a || y22+v>b )
        error('Window is out of image!') ;
    end
    % Build Lucas-Kanade equation (compute A, b)
    [S1,S2] = size(Ix1);
    Ix1 = Ix1(1+uu:S1-uu,1+vv:S2-vv);
    Iy1 = Iy1(1+uu:S1-uu,1+vv:S2-vv);
    A = [Ix1, Iy1];
    B = -Inte1(1+uu:S1-uu,1+vv:S2-vv)+Inte(1+uu:S1-uu,1+vv:S2-vv);
    [B1,B2] = size(B);
    % Compute delta_p
    ATA = (A')*A;
    %ATA=[sum(Ix1*Ix1','all'), sum(Ix1*Iy1','all');sum(Ix1*Iy1','all'),sum(Iy1*Iy1','all')];
    %sumA= ATA;
    [M,N] = size(ATA);
    sumA = [sum(ATA(1:Q,1:Q),'all'),sum(ATA(1:Q,Q+1:M),'all');sum(ATA(Q+1:M,1:9),'all'),sum(ATA(Q+1:M,Q+1:M),'all')];
    if isequal(sumA,[0,0;0,0])==1
        tf = randi([0,1],[2,1]);
        if tf(1) == 0
            u = -1*randi([1,8]);
        else
            u = randi([1,5]);
        end
        if tf(2) == 0
            v = -1*randi([1,8]);
        else
            v = randi([1,5]);
        end
        i=0;
        disp("Inv A is inf");
        continue
    end
    sumA = inv(sumA);
    ATB=[trace(Ix1'*B);trace(Iy1'*B)];
    %ATB=zeros(2,1);
    %for rr=1:B1
    %    for kk=1:B2
    %        ATB(1)=ATB(1)+Ix1(rr,kk)*B(rr,kk);
    %        ATB(2)=ATB(2)+Iy1(rr,kk)*B(rr,kk);
    %    end
    %end
    %disp(sumA);
    %disp(ATB);
    delta_p = sumA*ATB;
    %disp(delta_p);
    %disp("//delta norm///");
    % Update u, v
    u = u+delta_p(1);
    v = v+delta_p(2);
    % Compute the norm of delta_p
    delta_p_norm = delta_p(1)^2+delta_p(2)^2;
    if(i>=6&&delta_p_norm>0.001)
         tf = randi([0,1],[2,1]);
        if tf(1) == 0
            u = -1*randi([1,8]);
        else
            u = randi([1,5]);
        end
        if tf(2) == 0
            v = -1*randi([1,8]);
        else
            v = randi([1,5]);
        end
        i=0;
    end
    i = i+1;
end
disp(delta_p_norm);

% ------------------------ end of your code

end