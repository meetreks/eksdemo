package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func EC2CreateTags(resources []string, tags map[string]string) error {
	sess := GetSession()
	svc := ec2.New(sess)

	_, err := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: aws.StringSlice(resources),
		Tags:      createEC2Tags(tags),
	})

	if err != nil {
		return FormatError(err)
	}
	return nil
}

func EC2DescribeNetworkInterfaces(id, vpcId, description, instanceId, ip, securityGroupId string) ([]*ec2.NetworkInterface, error) {
	sess := GetSession()
	svc := ec2.New(sess)

	filters := []*ec2.Filter{}
	networkInterfaces := []*ec2.NetworkInterface{}
	pageNum := 0

	if id != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("network-interface-id"),
			Values: aws.StringSlice([]string{id}),
		})
	}

	if description != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("description"),
			Values: aws.StringSlice([]string{description}),
		})
	}

	if instanceId != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("attachment.instance-id"),
			Values: aws.StringSlice([]string{instanceId}),
		})
	}

	if ip != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("addresses.private-ip-address"),
			Values: aws.StringSlice([]string{ip}),
		})
	}

	if securityGroupId != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("group-id"),
			Values: aws.StringSlice([]string{securityGroupId}),
		})
	}

	if vpcId != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("vpc-id"),
			Values: aws.StringSlice([]string{vpcId}),
		})
	}

	input := &ec2.DescribeNetworkInterfacesInput{}

	if len(filters) > 0 {
		input.Filters = filters
	}

	err := svc.DescribeNetworkInterfacesPages(input,
		func(page *ec2.DescribeNetworkInterfacesOutput, lastPage bool) bool {
			pageNum++
			networkInterfaces = append(networkInterfaces, page.NetworkInterfaces...)
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, err
	}

	return networkInterfaces, nil
}

func EC2DescribeSecurityGroupRules(id, securityGroupId string) ([]*ec2.SecurityGroupRule, error) {
	sess := GetSession()
	svc := ec2.New(sess)

	filters := []*ec2.Filter{}
	securityGroupRules := []*ec2.SecurityGroupRule{}
	pageNum := 0

	if id != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("security-group-rule-id"),
			Values: aws.StringSlice([]string{id}),
		})
	}

	if securityGroupId != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("group-id"),
			Values: aws.StringSlice([]string{securityGroupId}),
		})
	}

	input := &ec2.DescribeSecurityGroupRulesInput{}

	if len(filters) > 0 {
		input.Filters = filters
	}

	err := svc.DescribeSecurityGroupRulesPages(input,
		func(page *ec2.DescribeSecurityGroupRulesOutput, lastPage bool) bool {
			pageNum++
			securityGroupRules = append(securityGroupRules, page.SecurityGroupRules...)
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, err
	}

	return securityGroupRules, nil
}

func EC2DescribeSecurityGroups(id, vpcId string, ids []string) ([]*ec2.SecurityGroup, error) {
	sess := GetSession()
	svc := ec2.New(sess)

	filters := []*ec2.Filter{}
	securityGroups := []*ec2.SecurityGroup{}
	pageNum := 0

	if id != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("group-id"),
			Values: aws.StringSlice([]string{id}),
		})
	}

	if vpcId != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("vpc-id"),
			Values: aws.StringSlice([]string{vpcId}),
		})
	}

	input := &ec2.DescribeSecurityGroupsInput{}

	if len(filters) > 0 {
		input.Filters = filters
	}

	if len(ids) > 0 {
		input.GroupIds = aws.StringSlice(ids)
	}

	err := svc.DescribeSecurityGroupsPages(input,
		func(page *ec2.DescribeSecurityGroupsOutput, lastPage bool) bool {
			pageNum++
			securityGroups = append(securityGroups, page.SecurityGroups...)
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, err
	}

	return securityGroups, nil
}

func EC2DescribeSubnets(id, vpcId string) ([]*ec2.Subnet, error) {
	sess := GetSession()
	svc := ec2.New(sess)

	filters := []*ec2.Filter{}
	subnets := []*ec2.Subnet{}
	pageNum := 0

	if id != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("subnet-id"),
			Values: aws.StringSlice([]string{id}),
		})
	}

	if vpcId != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("vpc-id"),
			Values: aws.StringSlice([]string{vpcId}),
		})
	}

	input := &ec2.DescribeSubnetsInput{}

	if len(filters) > 0 {
		input.Filters = filters
	}

	err := svc.DescribeSubnetsPages(input,
		func(page *ec2.DescribeSubnetsOutput, lastPage bool) bool {
			pageNum++
			subnets = append(subnets, page.Subnets...)
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, err
	}

	return subnets, nil
}

func EC2DescribeTags(resources, tagsFilter []string) (*ec2.DescribeTagsOutput, error) {
	sess := GetSession()
	svc := ec2.New(sess)

	filters := []*ec2.Filter{
		{
			Name:   aws.String("resource-id"),
			Values: aws.StringSlice(resources),
		},
	}

	if len(tagsFilter) > 0 {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("key"),
			Values: aws.StringSlice(tagsFilter),
		})
	}

	result, err := svc.DescribeTags(&ec2.DescribeTagsInput{
		Filters: filters,
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func EC2DescribeVpcs(id, vpcId string) ([]*ec2.Vpc, error) {
	sess := GetSession()
	svc := ec2.New(sess)

	filters := []*ec2.Filter{}
	vpcs := []*ec2.Vpc{}
	pageNum := 0

	if id != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("vpc-id"),
			Values: aws.StringSlice([]string{id}),
		})
	}

	if vpcId != "" {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String("vpc-id"),
			Values: aws.StringSlice([]string{vpcId}),
		})
	}

	input := &ec2.DescribeVpcsInput{}

	if len(filters) > 0 {
		input.Filters = filters
	}

	err := svc.DescribeVpcsPages(input,
		func(page *ec2.DescribeVpcsOutput, lastPage bool) bool {
			pageNum++
			vpcs = append(vpcs, page.Vpcs...)
			return pageNum <= maxPages
		},
	)

	if err != nil {
		return nil, err
	}

	return vpcs, nil
}

func createEC2Tags(tags map[string]string) (ec2tags []*ec2.Tag) {
	for k, v := range tags {
		ec2tags = append(ec2tags, &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return
}
