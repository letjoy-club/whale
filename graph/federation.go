// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graph

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/99designs/gqlgen/plugin/federation/fedruntime"
)

var (
	ErrUnknownType  = errors.New("unknown type")
	ErrTypeNotFound = errors.New("type not found")
)

func (ec *executionContext) __resolve__service(ctx context.Context) (fedruntime.Service, error) {
	if ec.DisableIntrospection {
		return fedruntime.Service{}, errors.New("federated introspection disabled")
	}

	var sdl []string

	for _, src := range sources {
		if src.BuiltIn {
			continue
		}
		sdl = append(sdl, src.Input)
	}

	return fedruntime.Service{
		SDL: strings.Join(sdl, "\n"),
	}, nil
}

func (ec *executionContext) __resolve_entities(ctx context.Context, representations []map[string]interface{}) []fedruntime.Entity {
	list := make([]fedruntime.Entity, len(representations))

	repsMap := map[string]struct {
		i []int
		r []map[string]interface{}
	}{}

	// We group entities by typename so that we can parallelize their resolution.
	// This is particularly helpful when there are entity groups in multi mode.
	buildRepresentationGroups := func(reps []map[string]interface{}) {
		for i, rep := range reps {
			typeName, ok := rep["__typename"].(string)
			if !ok {
				// If there is no __typename, we just skip the representation;
				// we just won't be resolving these unknown types.
				ec.Error(ctx, errors.New("__typename must be an existing string"))
				continue
			}

			_r := repsMap[typeName]
			_r.i = append(_r.i, i)
			_r.r = append(_r.r, rep)
			repsMap[typeName] = _r
		}
	}

	isMulti := func(typeName string) bool {
		switch typeName {
		default:
			return false
		}
	}

	resolveEntity := func(ctx context.Context, typeName string, rep map[string]interface{}, idx []int, i int) (err error) {
		// we need to do our own panic handling, because we may be called in a
		// goroutine, where the usual panic handling can't catch us
		defer func() {
			if r := recover(); r != nil {
				err = ec.Recover(ctx, r)
			}
		}()

		switch typeName {
		case "LevelRights":
			resolverName, err := entityResolverNameForLevelRights(ctx, rep)
			if err != nil {
				return fmt.Errorf(`finding resolver for Entity "LevelRights": %w`, err)
			}
			switch resolverName {

			case "findLevelRightsByLevel":
				id0, err := ec.unmarshalNInt2int(ctx, rep["level"])
				if err != nil {
					return fmt.Errorf(`unmarshalling param 0 for findLevelRightsByLevel(): %w`, err)
				}
				entity, err := ec.resolvers.Entity().FindLevelRightsByLevel(ctx, id0)
				if err != nil {
					return fmt.Errorf(`resolving Entity "LevelRights": %w`, err)
				}

				list[idx[i]] = entity
				return nil
			}
		case "Matching":
			resolverName, err := entityResolverNameForMatching(ctx, rep)
			if err != nil {
				return fmt.Errorf(`finding resolver for Entity "Matching": %w`, err)
			}
			switch resolverName {

			case "findMatchingByID":
				id0, err := ec.unmarshalNString2string(ctx, rep["id"])
				if err != nil {
					return fmt.Errorf(`unmarshalling param 0 for findMatchingByID(): %w`, err)
				}
				entity, err := ec.resolvers.Entity().FindMatchingByID(ctx, id0)
				if err != nil {
					return fmt.Errorf(`resolving Entity "Matching": %w`, err)
				}

				list[idx[i]] = entity
				return nil
			}
		case "MatchingQuota":
			resolverName, err := entityResolverNameForMatchingQuota(ctx, rep)
			if err != nil {
				return fmt.Errorf(`finding resolver for Entity "MatchingQuota": %w`, err)
			}
			switch resolverName {

			case "findMatchingQuotaByUserID":
				id0, err := ec.unmarshalNString2string(ctx, rep["userId"])
				if err != nil {
					return fmt.Errorf(`unmarshalling param 0 for findMatchingQuotaByUserID(): %w`, err)
				}
				entity, err := ec.resolvers.Entity().FindMatchingQuotaByUserID(ctx, id0)
				if err != nil {
					return fmt.Errorf(`resolving Entity "MatchingQuota": %w`, err)
				}

				list[idx[i]] = entity
				return nil
			}
		case "Topic":
			resolverName, err := entityResolverNameForTopic(ctx, rep)
			if err != nil {
				return fmt.Errorf(`finding resolver for Entity "Topic": %w`, err)
			}
			switch resolverName {

			case "findTopicByID":
				id0, err := ec.unmarshalNString2string(ctx, rep["id"])
				if err != nil {
					return fmt.Errorf(`unmarshalling param 0 for findTopicByID(): %w`, err)
				}
				entity, err := ec.resolvers.Entity().FindTopicByID(ctx, id0)
				if err != nil {
					return fmt.Errorf(`resolving Entity "Topic": %w`, err)
				}

				list[idx[i]] = entity
				return nil
			}
		case "User":
			resolverName, err := entityResolverNameForUser(ctx, rep)
			if err != nil {
				return fmt.Errorf(`finding resolver for Entity "User": %w`, err)
			}
			switch resolverName {

			case "findUserByID":
				id0, err := ec.unmarshalNString2string(ctx, rep["id"])
				if err != nil {
					return fmt.Errorf(`unmarshalling param 0 for findUserByID(): %w`, err)
				}
				entity, err := ec.resolvers.Entity().FindUserByID(ctx, id0)
				if err != nil {
					return fmt.Errorf(`resolving Entity "User": %w`, err)
				}

				list[idx[i]] = entity
				return nil
			}

		}
		return fmt.Errorf("%w: %s", ErrUnknownType, typeName)
	}

	resolveManyEntities := func(ctx context.Context, typeName string, reps []map[string]interface{}, idx []int) (err error) {
		// we need to do our own panic handling, because we may be called in a
		// goroutine, where the usual panic handling can't catch us
		defer func() {
			if r := recover(); r != nil {
				err = ec.Recover(ctx, r)
			}
		}()

		switch typeName {

		default:
			return errors.New("unknown type: " + typeName)
		}
	}

	resolveEntityGroup := func(typeName string, reps []map[string]interface{}, idx []int) {
		if isMulti(typeName) {
			err := resolveManyEntities(ctx, typeName, reps, idx)
			if err != nil {
				ec.Error(ctx, err)
			}
		} else {
			// if there are multiple entities to resolve, parallelize (similar to
			// graphql.FieldSet.Dispatch)
			var e sync.WaitGroup
			e.Add(len(reps))
			for i, rep := range reps {
				i, rep := i, rep
				go func(i int, rep map[string]interface{}) {
					err := resolveEntity(ctx, typeName, rep, idx, i)
					if err != nil {
						ec.Error(ctx, err)
					}
					e.Done()
				}(i, rep)
			}
			e.Wait()
		}
	}
	buildRepresentationGroups(representations)

	switch len(repsMap) {
	case 0:
		return list
	case 1:
		for typeName, reps := range repsMap {
			resolveEntityGroup(typeName, reps.r, reps.i)
		}
		return list
	default:
		var g sync.WaitGroup
		g.Add(len(repsMap))
		for typeName, reps := range repsMap {
			go func(typeName string, reps []map[string]interface{}, idx []int) {
				resolveEntityGroup(typeName, reps, idx)
				g.Done()
			}(typeName, reps.r, reps.i)
		}
		g.Wait()
		return list
	}
}

func entityResolverNameForLevelRights(ctx context.Context, rep map[string]interface{}) (string, error) {
	for {
		var (
			m   map[string]interface{}
			val interface{}
			ok  bool
		)
		_ = val
		m = rep
		if _, ok = m["level"]; !ok {
			break
		}
		return "findLevelRightsByLevel", nil
	}
	return "", fmt.Errorf("%w for LevelRights", ErrTypeNotFound)
}

func entityResolverNameForMatching(ctx context.Context, rep map[string]interface{}) (string, error) {
	for {
		var (
			m   map[string]interface{}
			val interface{}
			ok  bool
		)
		_ = val
		m = rep
		if _, ok = m["id"]; !ok {
			break
		}
		return "findMatchingByID", nil
	}
	return "", fmt.Errorf("%w for Matching", ErrTypeNotFound)
}

func entityResolverNameForMatchingQuota(ctx context.Context, rep map[string]interface{}) (string, error) {
	for {
		var (
			m   map[string]interface{}
			val interface{}
			ok  bool
		)
		_ = val
		m = rep
		if _, ok = m["userId"]; !ok {
			break
		}
		return "findMatchingQuotaByUserID", nil
	}
	return "", fmt.Errorf("%w for MatchingQuota", ErrTypeNotFound)
}

func entityResolverNameForTopic(ctx context.Context, rep map[string]interface{}) (string, error) {
	for {
		var (
			m   map[string]interface{}
			val interface{}
			ok  bool
		)
		_ = val
		m = rep
		if _, ok = m["id"]; !ok {
			break
		}
		return "findTopicByID", nil
	}
	return "", fmt.Errorf("%w for Topic", ErrTypeNotFound)
}

func entityResolverNameForUser(ctx context.Context, rep map[string]interface{}) (string, error) {
	for {
		var (
			m   map[string]interface{}
			val interface{}
			ok  bool
		)
		_ = val
		m = rep
		if _, ok = m["id"]; !ok {
			break
		}
		return "findUserByID", nil
	}
	return "", fmt.Errorf("%w for User", ErrTypeNotFound)
}
