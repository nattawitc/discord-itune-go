set l to (((path to desktop) as text) & "cover:")


tell application "Music"
	repeat with t in tracks of current playlist
		set rawdata to raw data of first artwork of t
		set a to album of t
		set a2 to my replace_chars(a, ":", ";")
		set l2 to (l & a2 & ".jpg")
		
		set outfile to open for access file l2 with write permission
		write rawdata to outfile
		close access outfile
	end repeat
end tell


on replace_chars(this_text, search_string, replacement_string)
	set AppleScript's text item delimiters to the search_string
	set the item_list to every text item of this_text
	set AppleScript's text item delimiters to the replacement_string
	set this_text to the item_list as string
	set AppleScript's text item delimiters to ""
	return this_text
end replace_chars