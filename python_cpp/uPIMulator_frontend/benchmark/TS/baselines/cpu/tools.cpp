#include <stdio.h>
#include <stdlib.h>
#include <math.h>
//#include <ctime>
#include <iostream>
#include <fstream>
#include <iomanip>
#include <limits>
#include <vector>
#include <algorithm>
#include <string>
#include <sstream>

#include "mprofile.h"


int loadTimeSeriesFromFile (std::string infilename, std::vector<DTYPE> &A, int &timeSeriesLength)
{
  std::fstream timeSeriesFile(std::string(PATH_TIME_SERIES) + infilename, std::ios_base::in);
	double tempval;
	timeSeriesLength = 0;
	while (timeSeriesFile >> tempval)
	{
		A.push_back(tempval);
		timeSeriesLength++;
	}
	timeSeriesFile.close();

  return 0;
}


int saveProfileToFile(std::string outfilename, DTYPE * profile, int * profileIndex, int timeSeriesLength, int windowSize)
{
  std::string preoutfilename = std::string(PATH_RESULTS) + outfilename;

  std::fstream preprofileOutFile(preoutfilename.c_str(), std::ios_base::out);

  // Write PreSCRIMP Matrix Profile and Matrix Profile Index to file.
  for (int i = 0; i < timeSeriesLength - windowSize + 1; i++)
  {
    preprofileOutFile << std::setprecision(std::numeric_limits<DTYPE>::digits10 + 2) << sqrt(abs(profile[i])) << " " << std::setprecision(std::numeric_limits<int>::digits10 + 1) << profileIndex[i] << std::endl;
  }

  preprofileOutFile.close();

  return 0;
}
