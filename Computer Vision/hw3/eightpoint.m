%*************************************************************************%
% Function:    eightpoint                                                 % 
% Description: Estimate the fundamental matrix F using N given            %
%              correspondences (N>=8)                                     %  
%                                                                         %
%              Input:  X and Y - each N × 2 matrices with coordinates     % 
%                                that constitute correspondences with     % 
%                                the first and second image respectively  %
%                                                                         %
%                      M - a scalar used to scale F which is equal to the %
%                          largest image dimension                        %
%                                                                         %
%              Output: F - the 3 x 3 fundamental matrix                   %
%*************************************************************************%

function [F] = eightpoint(X, Y, M)

% Check equal size
if ( (size(X,1) ~= size(Y,1)) || (size(X,2) ~= size(Y,2)) )
    error('X and Y must have the same size.') ;
end

% Need at least 8 point pairs
N = size(X, 1) ;
if (N < 8)
    error('At least 8 point pairs are required to compute F') ;
end

% ------- YOUR CODE HERE
Xnor = X/M;
Ynor = Y/M;
W = ones(N,9);
W(:,1)=Xnor(:,1).*Ynor(:,1);
W(:,2)=Xnor(:,2).*Ynor(:,1);
W(:,3)=Ynor(:,1);
W(:,4)=Xnor(:,1).*Ynor(:,2);
W(:,5)=Xnor(:,2).*Ynor(:,2);
W(:,6)=Ynor(:,2);
W(:,7)=Xnor(:,1);
W(:,8)=Xnor(:,2);
[U,S,V] = svd(W);
F_est = reshape(V(:,9),3,3)';
[u,s,v] = svd(F_est);
F_nor = u*diag([s(1) s(5) 0])*(v');
dia = ones(3,1)/M;
dia(3) = 1;
T = diag(dia);
F = (T')*F_nor*T;
% ------- END OF YOUR CODE

end
