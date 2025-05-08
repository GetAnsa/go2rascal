BEGIN {
  goTokenListLen=0;
  rascalTokenListLen=0;
  delete goTokenList;
  delete rascalTokenList;
  goType="";
  tokenType="";
  print "package main"
  print "import ("
  print "\t\"fmt\""
#  print "\t\"go/ast\""
  print "\t\"go/token\""
  print ")"
  print ""
}

NF == 1 && tokenType != "" && goType != "" {
#  print "reset line " $0
  print "func " tokenType "ToRascal(node " goType ") string {"
  print "\tswitch node {"
  i=0;
  for (idx in goTokenList) {
    print "\tcase " goTokenList[idx] ":"
    print "\t\treturn " rascalTokenList[idx]
  }
  print "\tdefault:"
  print "\t\tpanic(fmt.Sprintf(\"unknown" toupper(substr(tokenType, 1, 1)) substr(tokenType, 2) "(\\\"%s\\\"), node.String()))"
  print "\t}"
  print "}"
  print ""

  print "func rascal" toupper(substr(tokenType, 1, 1)) substr(tokenType, 2) "ToGo(rascal" toupper(substr(tokenType, 1, 1)) substr(tokenType, 2) " string) " goType " {"

  print "\tswitch rascal" toupper(substr(tokenType, 1, 1)) substr(tokenType, 2) "{"
  for (idx in goTokenList) {
    print "\tcase " rascalTokenList[idx] ":"
    print "\t\treturn " goTokenList[idx]
  }

  print "\tdefault:"
  print "\t\tpanic(fmt.Sprintf(\"unknown" toupper(substr(tokenType, 1, 1)) substr(tokenType, 2) "(\\\"%s\\\"), rascal" toupper(substr(tokenType, 1, 1)) substr(tokenType, 2) "))"
  print "\t}"
  print "}"
  print ""
  goTokenListLen=0;
  rascalTokenListLen=0;
  delete goTokenList;
  delete rascalTokenList;
  goType="";
  tokenType="";
}

NF == 1 && goType == "" && tokenType != "" {
  goType=$1;
#  print "goType " goType
}

NF == 1 && tokenType == "" {
  tokenType=$1;
#  print "tokenType " tokenType
}

# this will get assigned twice but I don't care
NF == 2 {
  goTokenList[goTokenListLen]=$1;
  rascalTokenList[rascalTokenListLen]=$2;
  goTokenListLen++;
  rascalTokenListLen++;
}

END {
  print "func " tokenType "ToRascal(node " goType ") string {"
  print "\tswitch node {"
  i=0;
  for (idx in goTokenList) {
    print "\tcase " goTokenList[idx] ":"
    print "\t\treturn " rascalTokenList[idx]
  }
  print "\tdefault:"
  print "\t\tpanic(fmt.Sprintf(\"unknown" toupper(substr(tokenType, 1, 1)) substr(tokenType, 2) "(\\\"%s\\\"), node.String()))"
  print "\t}"
  print "}"
  print ""

  print "func rascal" toupper(substr(tokenType, 1, 1)) substr(tokenType, 2) "ToGo(rascal" toupper(substr(tokenType, 1, 1)) substr(tokenType, 2) " string) " goType " {"

  print "\tswitch rascal" toupper(substr(tokenType, 1, 1)) substr(tokenType, 2) "{"
  for (idx in goTokenList) {
    print "\tcase " rascalTokenList[idx] ":"
    print "\t\treturn " goTokenList[idx]
  }

  print "\tdefault:"
  print "\t\tpanic(fmt.Sprintf(\"unknown" toupper(substr(tokenType, 1, 1)) substr(tokenType, 2) "(\\\"%s\\\"), rascal" toupper(substr(tokenType, 1, 1)) substr(tokenType, 2) "))"
  print "\t}"
  print "}"
}


