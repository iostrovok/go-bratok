#!/usr/bin/perl

use strict;

while (<STDIN>) {
    chomp $_;
    print STDERR $_, "\n";
    print STDOUT reverse(split '', $_), "\n";
    last;
}
