BEGIN {
  goToken="";
  rascalToken="";
}


/case *token.*:/ {
  rascalToken="";
  gsub(":","");
  goToken=$2;
}

/return *"/ {
  rascalToken=$2;

  print "\tcase " rascalToken ":"
  print "\t\treturn " goToken
}
