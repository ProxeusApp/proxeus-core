package test

type role int

type permission []byte

type groupAndOthers struct {
	//Allowed to modify: @Owner only!
	Group role `json:"group,omitempty"`
	//Rights pattern:       	    					group others
	//		                     	 	 				 --     --
	//default value is:                	    				----
	//example for group and others with read perm:	    	r-r-
	//Allowed to modify: @Owner only!
	Rights permission `json:"rights,omitempty"`
}

type permissions struct {
	//Allowed to modify: @Owner only!
	Owner string `json:"owner,omitempty"`
	//Grant can be modified by the owner only and is an optional field to whitelist user', s directly
	//Allowed to modify: everyone with write rights!
	Grant map[string]permission `json:"grant,omitempty"`

	GroupAndOthers groupAndOthers `json:"groupAndOthers,omitempty"`
	//Accessible by everyone who has the ID
	//Allowed to modify: everyone with write rights!
	PublicByID permission `json:"publicByID,omitempty"`

	//Execute only! If read or write not set.
	Published bool `json:"published"`
}
