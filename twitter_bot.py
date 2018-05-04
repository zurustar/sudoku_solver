#!/usr/bin/env python

import tweepy
import json
import sys
import re
from subprocess import Popen, PIPE


class SudokuStreamListener(tweepy.StreamListener):

	def set_username(self, username):
		self.username = username

	def set_api(self, api):
		self.api = api

	def on_status(self, status):
		#print(status.text, status.author)
		q = status.text.replace(self.username, '')
		q = re.sub('[^0-9]', '', q)
		#print("try to solve", q, "from", status.author.screen_name)
		p = Popen(['./sudoku_solver', q], stdout=PIPE, stderr=PIPE)
		out, err = p.communicate()
		out, err = out.decode('utf-8'), err.decode('utf-8') 
		self.api.update_status("@" + status.author.screen_name + "\n" + out + err,
                               in_reply_to_status_id=status.id)
		print(out, err)


def main(username, consumer_key, consumer_secret,
         access_token, access_token_secret):
	auth = tweepy.OAuthHandler(consumer_key, consumer_secret)
	auth.set_access_token(access_token, access_token_secret)
	api = tweepy.API(auth)
	listener=SudokuStreamListener()
	listener.set_username(username)
	listener.set_api(api)
	stream = tweepy.Stream(auth=api.auth, listener=listener)
	stream.filter(track=[username])

if __name__ == '__main__':
	fp = open(sys.argv[1], 'r')
	conf = json.load(fp)
	main(conf['username'],
         conf['consumer_key'], conf['consumer_secret'],
		 conf['access_token'], conf['access_token_secret'])
