#!/usr/bin/env bash

# <xbar.title>Hello World</xbar.title>
# <xbar.version>v1.0</xbar.version>
# <xbar.author>Jacob LeGrone</xbar.author>
# <xbar.author.github>jlegrone</xbar.author.github>
# <xbar.desc>Greet the current user</xbar.desc>
# <xbar.abouturl>https://github.com/jlegrone/xbargo</xbar.abouturl>
# <xbar.image>https://example.com/my-xbar-plugin-screenshot.png</xbar.image>
# <xbar.dependencies>github.com/jlegrone/xbargo/examples/helloworld</xbar.dependencies>
# <xbar.var>string(GREETING="hello"): The salutation to use</xbar.var>

# Set defaults again just in case they were deleted in plugin settings:
export GREETING="${GREETING:-"hello"}";

# Render menu items
./helloworld
