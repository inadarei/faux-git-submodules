#!/usr/bin/env python3

import json
import os
import sys

def getModules():

	try: 
		with open('fgs.json') as f:
			config = json.load(f)
	except FileNotFoundError:
		displayHelp(True)
		sys.exit(1)


	# Add the default branch, if not explicitly indicated in the config.
	config = {
		k: {**config[k], 'branch': config[k].get('branch', 'main')}
		for k in config		
	}
	return config

def checkPath(path):
	if not os.path.exists(path):
		return "missing"
	if os.path.isdir(path):
		return "folder"
	return "notFolder"

def execGitCommands(modules):
	"""
	Turn config into appropriate git commands. If a destination already
	exists then we do git pull, if not: git clone
	"""

	for module in modules:
		branch = modules[module]['branch']
		path = module
		url = modules[module]['url']
		
		whatIsIt = checkPath(module)
		if whatIsIt == "missing":
			cmd = f"git clone -b {branch} {url} {path}"
			# print(cmd)
			os.system(cmd)
		if whatIsIt == "folder":
			cmd = f"cd {module} && git pull"
			# print(cmd)
			os.system(cmd)
		if whatIsIt == "notFolder":
			print ("Cannot create '{module}' folder because path already exists")

def addIfMissing(ignorefile, entry):
	if not os.path.exists(ignorefile):
		os.system(f"touch {ignorefile}")

	with open(ignorefile,mode='r') as file:
		contents = file.read()

	if not entry in contents:
		with open(ignorefile, 'a') as fo:
			fo.write(f"{entry}\n")	
			print (f"Added {entry} to {ignorefile}");
	#else:
	#	print (f"ignorefile has {entry}")

def displayHelp(missingConfig=False):
	iam = os.path.basename(__file__)

	if missingConfig == True:
		print (f"Missing: fgs.json configuration file!")
		print("For an example fgs.json, see: https://github.com/inadarei/faux-git-submodules")
		print("")

	print ("Command utility to easily check-out faux git submodules.")
	print("Usage: ./"+f"{iam}")
	print ("For more information: https://github.com/inadarei/faux-git-submodules")

#-------------------------- Main Logic

modules = getModules();
execGitCommands(modules);

for reponame in modules:
	addIfMissing(".gitignore", reponame)

