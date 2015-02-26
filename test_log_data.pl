#!/usr/bin/perl

use strict;

my $s = "";
for (1 ... 110 ) {
	$s .= "$_: 00000000000000000000000000000--->STDOUT<--\n";
}

print $s;

my $e = "";
for (1 ... 110 ) {
	$e .= "$_: 00000000000000000000000000000--->ERROR+<--\n";
}

print STDERR $e; 

#print length($s), "\n";
#print length($e), "\n";
