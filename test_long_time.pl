#!/usr/bin/perl

my $duration = int shift;

print STDOUT "STDOUT BEFOR sleep: 0000000\n"; 
print STDERR "STDERR BEFOR sleep: EEEEEEE\n";

sleep($duration);

print STDOUT "STDOUT AFTER sleep: 0000000\n"; 
print STDERR "STDERR AFTER sleep: EEEEEEE\n";
