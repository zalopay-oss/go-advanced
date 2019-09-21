#!/bin/bash
gitbook build
gitbook epub . ./resource/advanced-go-book.epub
gitbook pdf . ./resource/advanced-go-book.pdf
gitbook mobi . ./resource/advanced-go-book.mobi